package handler

import (
	//"fmt"
	"time"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"net/url"
	"git.lianjia.com/lianjia-sysop/napi/variable"
	"git.lianjia.com/lianjia-sysop/napi/log"
	"git.lianjia.com/lianjia-sysop/napi/errors"
)

type OrchResponse struct {
	Jobid string
}

// Orch Orch handler
type OrchHandler struct {
	url       string
	createLoc string
	cancelLoc string
	statLoc   string
	teeLoc    string

	group     string
	token     string

	cb        string
	retry     int

	log       log.Log
}

var orchWorker *OrchHandler

// InitHandler inits handler
func InitOrchHandler(url, createLoc, cancelLoc, statLoc, teeLoc, group, token, cb string, log log.Log) {
	h := &OrchHandler{}
	h.url = url
	h.createLoc = createLoc
	h.cancelLoc = cancelLoc
	h.statLoc = statLoc
	h.teeLoc = teeLoc
	h.group = group
	h.token = token
	h.cb = cb
	h.retry = ORCH_RETRY_NUM
	h.log = log

	orchWorker = h
}

// Operate operates orch
func (oh *OrchHandler) Operate(op int, args string) (string, error) {
	var err error
	var floc string
	var resp *http.Response

	retry := 0

next:
	switch op {
	case ORCH_CREATE:
		oh.log.Debug("Orch create job args is %s", args)
		floc = "http://" + oh.url + oh.createLoc
		break
	case ORCH_CANCEL:
		oh.log.Debug("Orch cancel job args is %s", args)
		floc = "http://" + oh.url + oh.cancelLoc
		break
	case ORCH_READ:
		oh.log.Debug("Orch read job args is %s", args)
		floc = "http://" + oh.url + oh.statLoc
		break
	default: /* Should not reach here */
		oh.log.Error("Unknown operate code: ", op)
		return "", errors.InternalServerError
	}

	//val := &url.Values{}
	//val.Add("value", args)
	data := bytes.NewBufferString(args)

	oh.log.Debug("orch full url is %s",  floc)
	req, _ := http.NewRequest("POST", floc , data)

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set(ORCH_GROUP_HEADER, oh.group)
	req.Header.Set(ORCH_TOKEN_HEADER, oh.token)

	client := &http.Client{}
	resp, err = client.Do(req)

	if err != nil {
		oh.log.Error("Opereate job to orch failed: ", err)
		retry++
		if retry >= oh.retry {
			return "", errors.BadGatewayError
		}
		time.Sleep(time.Second)
		goto next
	}

	defer resp.Body.Close()

	if resp.StatusCode != variable.HTTP_OK {
		oh.log.Debug("Opereate http status error: %d", resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			oh.log.Error("Read operate result body failed: ", err)
			return "", errors.InternalServerError
		}
		oh.log.Debug("orch body is %s", body)
		//if (op != ADD) && (resp.StatusCode == variable.HTTP_NOT_FOUND) {
		//	return "", errors.NoContentError
		//}
		return "", errors.BadGatewayError
	}

	/* read operation */
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		oh.log.Error("Read operate result body failed: ", err)
		return "", errors.InternalServerError
	}

	eresp := &OrchResponse{}
	err = json.Unmarshal(body, &eresp)
	if err != nil {
		oh.log.Error("Unmarshl from response failed")
		return "", errors.InternalServerError
	}

	return eresp.Jobid, nil
}

// Create create job to orch
func OrchCreate(om *OrchMessage) (string, error) {
	om.Callback = orchWorker.cb
	om.Type = "exec"
	b, err := json.Marshal(om)
	if err != nil {
		orchWorker.log.Error("Orch create marshal failed")
		return "", err
	}
	args := string(b)

	orchWorker.log.Debug(args)

	return orchWorker.Operate(ORCH_CREATE, args)
}

// Cancel cancel job to orch
func OrchCancel(om *OrchMessage) (string, error) {
	b, err := json.Marshal(om)
	if err != nil {
		orchWorker.log.Error("Orch cancel marshal failed")
		return "", err
	}
	args := string(b)

	orchWorker.log.Debug(args)

	return orchWorker.Operate(ORCH_CANCEL, args)
}

// Read read job from orch
func OrchRead(om *OrchMessage) (string, error) {
	b, err := json.Marshal(om)
	if err != nil {
		orchWorker.log.Error("Orch read marshal failed")
		return "", err
	}
	args := string(b)

	orchWorker.log.Debug(args)

	return orchWorker.Operate(ORCH_READ, args)
}
