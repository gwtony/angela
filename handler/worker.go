package handler

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"github.com/gwtony/angela/ssh"
	"github.com/gwtony/gapi/log"
)

func RunJob(job OrchMessage, log log.Log) bool {

	if job.Callback != "" {
		//startJobCallback(jobObj)
		//TODO: call or not
	}

	defer ReportJobCallback(job)

	concurrency := MAX_CONCURRENT
	if job.Parallel <= 0 {
		concurrency = 1
	}

	g.InitTmpDir()
	tmpFile := g.Config().Tmpdir + "/" + jobObj.Jobid + ".tmp"

	defer func(debug bool, tmpf string) {
		if !debug {
			os.Remove(tmpf)
		}
	}(g.Config().Debug, tmpFile)

	werr := ioutil.WriteFile(tmpFile, []byte(job.Command), 0644)
	if werr != nil {
		//log.Println("Save JobCommand To File Fail!", tmpFile, werr)
		for _, node := range job.Nodes {
			updateJobNodeState(job.Jobid, node, g.JOB_FAIL_STR, -1)
			updateJobNodeOutput(job.Jobid, node, "", "Save job command to tmpfile fail", time.Now(), time.Now())
		}
		return false
	}

	//run command
	sem := make(chan bool, concurrency)
	for _, node := range job.Nodes {
		sem <- true

		go func(node string, jobid string, user string, timeout int) {
			defer func() { <-sem }()

			//if g.Config().Debug {
			//	log.Printf("concurrentcy pools length: %v, capability:%v\n", len(sem), cap(sem))
			//}

			//update node to running
			updateJobNodeState(jobid, node, g.JOB_RUNNING_STR, 0)

			//run command
			rc, isTout := executeCmd(jobid, node, user, tmpFile, timeout)

			execState := g.JOB_SUCC_STR
			if rc != 0 {
				execState = g.JOB_FAIL_STR
			}
			if !isTout {
				execState = g.JOB_TIMEOUT_STR
			}

			updateJobNodeState(job.Jobid, node, execState, rc)

		}(node, job.Jobid, job.User, job.Timeout)
	}

	//release concurrency
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return true
}

func executeCmd(jobid, node, user, tmpfile string, timeout int) (rcInt int, isTout bool) {
	if execuser == "" {
		return -1, true
	}
	ssh := &easyssh.MakeConfig{
		User:   "root",
		Server: node,
		Key:    g.Config().Sshkey,
		Port:   "22",
	}

	//if g.Config().Debug {
	//	log.Println("=======>ssh msg is:", node, g.Config().Sshkey)

	//	ssh = &easyssh.MakeConfig{
	//		User:     "root",
	//		Server:   node,
	//		Password: "123qwe",
	//		Port:     "22",
	//	}
	//}

	_, _, _, err := ssh.Run("mkdir /tmp/orch-run", 60)
	if err != nil {
		updateJobNodeOutput(jobid, node, "", derr.Error(), time.Now(), time.Now())
		return -1, true
	}

	//scp command file
	remoteFile := "/tmp/orch-run/" + jobid + ""
	err := ssh.Scp(tmpfile, remoteFile)
	if err != nil {
		updateJobNodeOutput(jobid, node, "", serr.Error(), time.Now(), time.Now())
		return -1, true
	}

	//get script shell，default is "bash -l", act as a login shell
	shellCmd := "bash -l "
	f, err := os.Open(tmpfile)
	if err != nil {
		updateJobNodeOutput(jobid, node, "", serr.Error(), time.Now(), time.Now())
		return -1, true
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

	//delete command file
	defer ssh.Run("rm -f "+remoteFile, 60)
	//kill command
	defer ssh.Run("ps auxf|grep '"+remoteFile+"'|grep -v grep|awk '{print $2}'|xargs kill ", 60)

	//run command
	start := time.Now()
	command := "su " + execuser + " -s /bin/bash -c 'source /etc/profile; " + shellCmd + "  " + remoteFile + "; echo $?'"
	if g.Config().Debug {
		log.Println("======> Command is: ", command)
	}
	stdout, stderr, isTout, err := ssh.Run(command, timeout)
	if err != nil {
		updateJobNodeOutput(jobid, node, "", err.Error(), start, start)
		return -1, true
	}
	if g.Config().Debug {
		log.Println("======> Command is: ", command, " , run command stdout is :", stdout, " ; stderr is :", stderr, " ; isTout is :", isTout, " ; err is :", err)
	}

	end := time.Now()

	if isTout {
		//get return code from stdout
		stdoutArr := strings.Split(strings.TrimSpace(stdout), "\n")
		rc := stdoutArr[len(stdoutArr)-1]
		if strings.TrimSpace(rc) == "" {
			rcInt = -1
		} else {
			rcInt = utils.ToInt(rc)
		}
	} else {
		rcInt = -1
	}
	log.Println("======> kill process command:", "ps auxf|grep '"+remoteFile+"'|grep -v grep|awk '{print $2}'|xargs kill ")

	updateJobNodeOutput(jobid, node, stdout, stderr, start, end)

	return rcInt, isTout
}

func startJobCallback(jobObj m.JJob) {
	if jobObj.Callback == "" {
		return
	}

	var reqObj m.JJobCallback
	reqObj.Jobid = jobObj.Jobid
	reqObj.State = make([]m.JNodeStatus, 0)
    reqObj.Meta = jobObj.Meta
	for _, n := range jobObj.Nodes {
		var nObj m.JNodeStatus
		nObj.Node = n
		nObj.State = g.JOB_RUNNING_STR
		reqObj.State = append(reqObj.State, nObj)
	}

	b, err := json.Marshal(reqObj)
	if g.Config().Debug {
		log.Printf("startJobCallback: %s; request json %s \n", jobObj.Callback, string(b))
	}
	if err != nil {
		if g.Config().Debug {
			log.Printf("startJobCallback: %v, json encode error! %v \n", jobObj, err)
		}
		return
	}

	req := httplib.Post(jobObj.Callback)
	req.Body(b)
	req.SetTimeout(10*time.Second, 10*time.Second)
	resp, err1 := req.Response()
	if err1 != nil {
		if g.Config().Debug {
			log.Printf("startJobCallback: %v, response error! %v \n", jobObj, err1)
		}
		return
	}

	status := resp.StatusCode
	if g.Config().Debug {
		log.Printf("startJobCallback: %v, response status: %v \n", jobObj, status)
	}

	return
}

func finishJobCallback(jobObj m.JJob) {
	if jobObj.Callback == "" {
		return
	}

	var reqObj m.JJobCallback
	reqObj.Jobid = jobObj.Jobid
    reqObj.Meta = jobObj.Meta

	jsList, _ := db.QueryJobStates(jobObj.Jobid)
	for _, n := range jsList {
		var nObj m.JNodeStatus
		nObj.Node = n.Node
		nObj.State = n.State
		nObj.Rc = n.Rc
		reqObj.State = append(reqObj.State, nObj)
	}

	jsOutputList, _ := db.QueryJobOutput(jobObj.Jobid)
	for _, n := range jsOutputList {
		var nObj m.JNodeOutput
		nObj.Node = n.Node
		nObj.Stdout = n.Stdout
		nObj.Stderr = n.Stderr
		nObj.Start = n.Start.Format("2006-01-02 15:04:05")
		nObj.End = n.End.Format("2006-01-02 15:04:05")
		nObj.Delta = n.Delta
		reqObj.Output = append(reqObj.Output, nObj)
	}

	b, err := json.Marshal(reqObj)
	if err != nil {
		if g.Config().Debug {
			log.Printf("finishJobCallback: %v, json encode error! %v \n", jobObj, err)
		}
		return
	}
	if g.Config().Debug {
		log.Printf("finishJobCallback: %s; request json %s \n", jobObj.Callback, string(b))
	}

	req := httplib.Post(jobObj.Callback)
	req.Body(b)
	req.SetTimeout(10*time.Second, 10*time.Second)
	resp, err1 := req.Response()
	if err1 != nil {
		if g.Config().Debug {
			log.Printf("finishJobCallback: %v, response error! %v \n", jobObj, err1)
		}
		return
	}

	status := resp.StatusCode
	if g.Config().Debug {
		log.Printf("finishJobCallback: %v, response status: %v \n", jobObj, status)
	}

	return
}

func updateJobNodeState(jobid, node, state string, rc int) bool {
	if g.Config().Debug {
		log.Println("----------> updateJobNodeState:", jobid, node, state, rc)
	}
	//update node state to running
	var jsObj m.Jobstate
	jsObj.Jobid = jobid
	jsObj.Node = node
	jsObj.State = state
	jsObj.Rc = rc
	db.UpdateJobState(jsObj)
	return true
}

func updateJobNodeOutput(jobid, node, stdout, stderr string, start, end time.Time) bool {
	if g.Config().Debug {
		log.Println("----------> updateJobNodeOutput:", jobid, node, stdout, stderr, start, end)
	}
	var jnOutput m.Joboutput
	jnOutput.Jobid = jobid
	jnOutput.Node = node
	jnOutput.Stdout = stdout
	jnOutput.Stderr = stderr
	jnOutput.Start = start
	jnOutput.End = end
	jnOutput.Delta = utils.ToInt(end.Unix() - start.Unix())
	db.InsertJobOutput(db.NewOrm(), jnOutput)
	return true
}
