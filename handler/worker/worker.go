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

//var MAX_CONCURRENT int = 8
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
		//startJobCallback(jobObj)
		//TODO: call or not
	}

	//TODO: config concurrency
	//concurrency := 8
	//if job.Parallel <= 0 {
	//	concurrency = 1
	//}

	//g.InitTmpDir()
	jobid := job.Meta.Runid
	tmpFile := variable.ORCH_TMP_DIR + "/" + jobid + ".tmp"

	//defer os.Remove(tmpf)

	err := ioutil.WriteFile(tmpFile, []byte(job.Command), 0644)
	if err != nil {
		log.Error("Job[%s] Save job command to file fail!", jobid, tmpFile, err)
		//for _, node := range job.Nodes {
		//	updateJobNodeState(job.Jobid, node, JOB_FAILED, -1)
		//	updateJobNodeOutput(job.Jobid, node, "", "Save job command to tmpfile fail", time.Now(), time.Now())
		//}
		ReportJobCallback(job, nil, log)
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

			//update node to running
			//updateJobNodeState(jobid, node, JOB_RUNNING, 0)

			//execute command
			err := execute(stateCh, jid, node, user, tmpFile, timeout, log)
			if err != nil {
				log.Error("Execute job[%s] failed:", jid, err)
			}

			//updateJobNodeState(job.Jobid, node, execState, rc)

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
			rs.Rc = -1 //Output = s.Error.Error()
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

	ReportJobCallback(job, rm, log)
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
		Port:   "22",
	}

	//if Debug {
	//	log.Println("=======>ssh msg is:", node, g.Config().Sshkey)

	//	ssh = &easyssh.MakeConfig{
	//		User:     "root",
	//		Server:   node,
	//		Password: "123qwe",
	//		Port:     "22",
	//	}
	//}

	_, _, _, err := session.Run("mkdir " + variable.ORCH_RUNTIME_DIR, 60)
	if err != nil {
		//updateJobNodeOutput(jobid, node, "", err.Error(), time.Now(), time.Now())
		state.Error = err
		return err
	}

	//scp command file
	remoteFile := variable.ORCH_RUNTIME_DIR + "/" + jobid + ""
	log.Debug("scp %s to %s", tmpfile, remoteFile)
	err = session.Scp(tmpfile, remoteFile)
	if err != nil {
		//updateJobNodeOutput(jobid, node, "", err.Error(), time.Now(), time.Now())
		log.Error("Scp script failed:", err)
		state.Error = err
		return err
	}

	//get script shell，default is "bash -l", act as a login shell
	shellCmd := "bash "
	//shellCmd := "bash -l "
	f, err := os.Open(tmpfile)
	if err != nil {
		//updateJobNodeOutput(jobid, node, "", err.Error(), time.Now(), time.Now())
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
		//updateJobNodeOutput(jobid, node, "", err.Error(), start, start)
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

	//updateJobNodeOutput(jobid, node, stdout, stderr, start, end)
	state.Node    = node
	state.Stdout  = stdout
	state.Stderr  = stderr
	state.Timeout = isTo
	state.Rc      = rc
	state.Error   = nil
	state.Start   = start.String() //TODO:
	state.End     = end.String()
	state.Delta   = uint64(end.Unix() - start.Unix())
	//fmt.Printf("state is", state)

	return nil
}

func ReportJobCallback(job msg.OrchMessage, rm *msg.ReportMessage, log log.Log) {
	log.Debug("job, rm:", job, rm)
}

//func startJobCallback(jobObj m.JJob, log log.Log) {
//	if jobObj.Callback == "" {
//		return
//	}
//
//	var reqObj m.JJobCallback
//	reqObj.Jobid = jobObj.Jobid
//	reqObj.State = make([]m.JNodeStatus, 0)
//    reqObj.Meta = jobObj.Meta
//	for _, n := range jobObj.Nodes {
//		var nObj m.JNodeStatus
//		nObj.Node = n
//		nObj.State = JOB_RUNNING
//		reqObj.State = append(reqObj.State, nObj)
//	}
//
//	b, err := json.Marshal(reqObj)
//	log.Debug("startJobCallback: %s; request json %s \n", jobObj.Callback, string(b))
//	if err != nil {
//		log.Debug("startJobCallback: %v, json encode error! %v \n", jobObj, err)
//		return
//	}
//
//	req := httplib.Post(jobObj.Callback)
//	req.Body(b)
//	req.SetTimeout(10*time.Second, 10*time.Second)
//	resp, err1 := req.Response()
//	if err1 != nil {
//		log.Debug("startJobCallback: %v, response error! %v \n", jobObj, err1)
//		return
//	}
//
//	status := resp.StatusCode
//	log.Debug("startJobCallback: %v, response status: %v \n", jobObj, status)
//
//	return
//}
//
//func finishJobCallback(jobObj m.JJob, log log.Log) {
//	if jobObj.Callback == "" {
//		return
//	}
//
//	var reqObj m.JJobCallback
//	reqObj.Jobid = jobObj.Jobid
//    reqObj.Meta = jobObj.Meta
//
//	jsList, _ := db.QueryJobStates(jobObj.Jobid)
//	for _, n := range jsList {
//		var nObj m.JNodeStatus
//		nObj.Node = n.Node
//		nObj.State = n.State
//		nObj.Rc = n.Rc
//		reqObj.State = append(reqObj.State, nObj)
//	}
//
//	jsOutputList, _ := db.QueryJobOutput(jobObj.Jobid)
//	for _, n := range jsOutputList {
//		var nObj m.JNodeOutput
//		nObj.Node = n.Node
//		nObj.Stdout = n.Stdout
//		nObj.Stderr = n.Stderr
//		nObj.Start = n.Start.Format("2006-01-02 15:04:05")
//		nObj.End = n.End.Format("2006-01-02 15:04:05")
//		nObj.Delta = n.Delta
//		reqObj.Output = append(reqObj.Output, nObj)
//	}
//
//	b, err := json.Marshal(reqObj)
//	if err != nil {
//		log.Debug("finishJobCallback: %v, json encode error! %v \n", jobObj, err)
//		return
//	}
//	log.Debug("finishJobCallback: %s; request json %s \n", jobObj.Callback, string(b))
//
//	req := httplib.Post(jobObj.Callback)
//	req.Body(b)
//	req.SetTimeout(10*time.Second, 10*time.Second)
//	resp, err1 := req.Response()
//	if err1 != nil {
//		log.Debug("finishJobCallback: %v, response error! %v \n", jobObj, err1)
//		return
//	}
//
//	status := resp.StatusCode
//
//	log.Debug("finishJobCallback: %v, response status: %v \n", jobObj, status)
//
//	return
//}

//func updateJobNodeState(jobid, node, state string, rc int, log log.Log) bool {
//	//update node state to running
//	var jsObj m.Jobstate
//
//	log.Debug("updateJobNodeState:", jobid, node, state, rc)
//
//	jsObj.Jobid = jobid
//	jsObj.Node = node
//	jsObj.State = state
//	jsObj.Rc = rc
//	db.UpdateJobState(jsObj)
//
//	return true
//}
//
//func updateJobNodeOutput(jobid, node, stdout, stderr string, start, end time.Time, log log.Log) bool {
//	var jnOutput m.Joboutput
//
//	log.Debug("updateJobNodeOutput:", jobid, node, stdout, stderr, start, end)
//
//	jnOutput.Jobid = jobid
//	jnOutput.Node = node
//	jnOutput.Stdout = stdout
//	jnOutput.Stderr = stderr
//	jnOutput.Start = start
//	jnOutput.End = end
//	jnOutput.Delta = utils.ToInt(end.Unix() - start.Unix())
//	db.InsertJobOutput(db.NewOrm(), jnOutput)
//
//	return true
//}
