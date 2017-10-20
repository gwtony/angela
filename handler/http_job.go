package handler
import (
	//"fmt"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/utils"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/errors"
)

type JobCreateHandler struct {
	mc  *MysqlContext
	log log.Log
}

type JobCancelHandler struct {
	mc  *MysqlContext
	log log.Log
}

func (h *JobCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var key string
	admin := false

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

	data := &CronMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}
	h.log.Info("Add record request: (%s) from client: %s", string(result), r.RemoteAddr)

	//TODO: check cron format
	if data.Cron == "" {
		h.log.Info("Not found cron")
		api.ReturnError(r, w, errors.Jerror("No cron in data"), errors.BadRequestError, h.log)
		return
	}

	if r.Header[ADMIN_TOKEN_HEADER] != nil {
		if (strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) != 0) {
			h.log.Info("Job add admin token invalid")
			api.ReturnError(r, w, errors.Jerror("Admin token unauthorized"), errors.UnauthorizedError, h.log)
			return
		}
		admin = true
	}

	jobid, _ := utils.NewUUID()

	im := &IdMessage{Jobid: data.Jobid}
	imv, _ := json.Marshal(im)

	api.ReturnResponse(r, w, string(imv), h.log)
}

func (h *JobCancelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	admin := false

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

	data := &CronOperateMessage{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
		return
	}

	h.log.Info("Cancel record request: (%s) from client: %s", string(result), r.RemoteAddr)

	api.ReturnResponse(r, w, "", h.log)
}
