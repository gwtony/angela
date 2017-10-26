package worker

import (
	"io"
	"os"
	//"fmt"
	"time"
	"sync"
	"errors"
	"bufio"
	"strings"
	"strconv"
	"io/ioutil"
	//"encoding/json"
	"github.com/gwtony/gapi/log"
	ssh "github.com/gwtony/angela/handler/ssh"
	"github.com/gwtony/angela/handler/msg"
	"github.com/gwtony/angela/handler/variable"
)

var keyPath string

func InitWorker(key string) error {
	keyPath = key

	fi, err := os.Lstat(variable.ORCH_TMP_DIR)
	if err != nil || fi.Mode().IsDir() == false {
		err = os.Mkdir(variable.ORCH_TMP_DIR, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunJob(job msg.OrchMessage, log log.Log) {
	if job.Callback != "" {
		//JobCallback(job, start, nil, log)
		//TODO: call or not
	}

	//TODO: config concurrency

	jobid := job.Meta.Runid
	tmpFile := variable.ORCH_TMP_DIR + "/" + jobid + ".tmp"

	//defer os.Remove(tmpf)

	err := ioutil.WriteFile(tmpFile, []byte(job.Command), 0644)
	if err != nil {
		log.Error("Job[%s] Save job command to file fail!", jobid, tmpFile, err)
		JobCallback(job, "failed", nil, log)
		return
	}

	//execute command
	//TODO: limit concurrency
	var wg sync.WaitGroup
	stateCh := make(chan msg.NodeState, len(job.Nodes))

	for _, node := range job.Nodes {
		wg.Add(1)

		go func(node, jid, user string, timeout int) {
			defer wg.Done()

			//execute command
			err := execute(stateCh, jid, node, user, tmpFile, timeout, log)
			if err != nil {
				log.Error("Execute job[%s] failed:", jid, err)
			}
		}(node, jobid, job.User, int(job.Timeout))
	}

	wg.Wait()

	rm := &msg.ReportMessage{}
	rm.Meta.Runid = jobid
	rm.State = make([]msg.ReportState, 0, len(job.Nodes))
	rm.Output = make([]msg.ReportOutput, 0, len(job.Nodes))

	for i := 0; i < len(job.Nodes); i++ {
		//TODO: report
		s := <-stateCh
		rs := msg.ReportState{}
		ro := msg.ReportOutput{}
		log.Debug("collect state")
		rs.Node = s.Node
		ro.Node = s.Node
		if s.Error != nil {
			rs.State = variable.JOB_FAILED
			rs.Rc = -1
			rs.Error = s.Error.Error()
			//TODO: fill ro
		} else {
			rs.State = variable.JOB_SUCCESS
			rs.Rc = s.Rc
			rs.Error = ""

			ro.Stdout = s.Stdout
			ro.Stderr = s.Stderr
			ro.Start = s.Start
			ro.Stop = s.End //TODO: change stop to end
			ro.Delta = uint64(s.Delta)
		}
		rm.State = append(rm.State, rs)
		rm.Output = append(rm.Output, ro)
	}

	JobCallback(job, "success", rm, log)
	return
}

func execute(ch chan msg.NodeState, jobid, node, user, tmpfile string, timeout int, log log.Log) error {
	state := msg.NodeState{Node: node}
	defer func() { ch <- state }()

	if user == "" {
		state.Error = errors.New("No execute user")
		return errors.New("No execute user")
	}
	log.Debug("key path is %s", keyPath)
	session := &ssh.MakeConfig{
		User:   "root",
		Server: node,
		Key:    keyPath,
		//Password: "111",
		Port:   "22",
	}

	_, _, _, err := session.Run("mkdir " + variable.ORCH_RUNTIME_DIR, 60)
	if err != nil {
		state.Error = err
		return err
	}

	//scp command file
	remoteFile := variable.ORCH_RUNTIME_DIR + "/" + jobid + ""
	log.Debug("scp %s to %s", tmpfile, remoteFile)
	err = session.Scp(tmpfile, remoteFile)
	if err != nil {
		log.Error("Scp script failed:", err)
		state.Error = err
		return err
	}

	//Act as a login shell or not
	shellCmd := "bash "
	//shellCmd := "bash -l "
	f, err := os.Open(tmpfile)
	if err != nil {
		state.Error = err
		return err
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		aline := strings.TrimSpace(string(a))
		if aline == "" {
			continue
		} else if strings.HasPrefix(aline, "#!") {
			shellCmd = strings.TrimPrefix(aline, "#!")
		} else {
			break
		}
	}

	//TODO: delete command file
	//defer session.Run("rm -f "+remoteFile, 60)
	//defer session.Run("ps auxf|grep '"+remoteFile+"'|grep -v grep|awk '{print $2}'|xargs kill ", 60)

	start := time.Now()
	//command := "su " + user + " -s /bin/bash -c 'source /etc/profile; " + shellCmd + " " + remoteFile + "; echo $?'"
	command := "su " + user + " -s /bin/bash -c '" + shellCmd + remoteFile + "; echo $?'"

	stdout, stderr, isTo, err := session.Run(command, timeout)
	if err != nil {
		log.Error("Command run failed:", err)
		state.Error = err
		return err
	}

	log.Debug("Stdout is: \"%s\"", stdout)
	if stderr != "" {
		log.Debug("Stderr is: \"%s\"", stderr)
	}
	if isTo {
		log.Debug("Command run timeout: ", isTo)
	}

	end := time.Now()

	rc := 0
	if isTo {
		//get return code from stdout
		stdoutArr := strings.Split(strings.TrimSpace(stdout), "\n")
		rcstr := stdoutArr[len(stdoutArr)-1]
		if strings.TrimSpace(rcstr) == "" {
			rc = -1
		} else {
			rc, _ = strconv.Atoi(rcstr)
		}
	} else {
		rc = -1
	}
	//log.Println("kill process command:", "ps auxf|grep '"+remoteFile+"'|grep -v grep|awk '{print $2}'|xargs kill ")
	log.Debug("to set state, rc: %d", rc)

	state.Node    = node
	state.Stdout  = stdout
	state.Stderr  = stderr
	state.Timeout = isTo
	state.Rc      = rc
	state.Error   = nil
	state.Start   = start.String() //TODO:
	state.End     = end.String()
	state.Delta   = uint64(end.Unix() - start.Unix())

	return nil
}

func JobCallback(job msg.OrchMessage, state string, rm *msg.ReportMessage, log log.Log) {
	log.Debug("state[%s], job, rm:", state, job, rm)
}

