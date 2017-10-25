package http_handler
import (
	//"fmt"
	//"time"
	//"strings"
	//"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	//"github.com/gwtony/gapi/utils"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/angela/handler/worker"
	"github.com/gwtony/angela/handler/msg"
	//"github.com/gwtony/angela/config"
	//"github.com/gwtony/angela/mysql"
)

type JobCreateHandler struct {
	log log.Log
}

type JobCancelHandler struct {
	log log.Log
}

func InitJobCreateHandler(log log.Log) *JobCreateHandler {
	return &JobCreateHandler{log: log}
}

func InitJobCancelHandler(log log.Log) *JobCancelHandler {
	return &JobCancelHandler{log: log}
}

func (h *JobCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//var key string
	//admin := false

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := msg.OrchMessage{}

	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Create job request: (%s) from client: %s", string(result), r.RemoteAddr)

	if len(data.Nodes) <= 0 {
		h.log.Info("Not found nodes")
		api.ReturnError(r, w, errors.Jerror("No nodes in data"), errors.BadRequestError, h.log)
		return
	}
	//TODO: check args

	//if r.Header[ADMIN_TOKEN_HEADER] != nil {
	//	if (strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) != 0) {
	//		h.log.Info("Job add admin token invalid")
	//		api.ReturnError(r, w, errors.Jerror("Admin token unauthorized"), errors.UnauthorizedError, h.log)
	//		return
	//	}
	//	admin = true
	//}

	//jobid, _ := utils.NewUUID()

	//im := &IdMessage{Jobid: data.Jobid}
	//imv, _ := json.Marshal(im)
	go worker.RunJob(data, h.log)

	//api.ReturnResponse(r, w, string(imv), h.log)
	api.ReturnResponse(r, w, "", h.log)
}

func (h *JobCancelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//admin := false

	if r.Method != "POST" {
		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
		return
	}
	r.Body.Close()

	data := &msg.OrchMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Cancel job request: (%s) from client: %s", string(result), r.RemoteAddr)

	api.ReturnResponse(r, w, "", h.log)
}
