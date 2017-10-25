package http_handler
//import (
//	"strings"
//	"io/ioutil"
//	"encoding/json"
//	"net/http"
//	"github.com/gwtony/gapi/log"
//	"github.com/gwtony/gapi/api"
//	"github.com/gwtony/gapi/errors"
//)
//
//type GroupAddHandler struct {
//	//mc  *MysqlContext
//	log log.Log
//}
//
//type GroupDeleteHandler struct {
//	//mc  *MysqlContext
//	log log.Log
//}
//
//type GroupReadHandler struct {
//	//mc  *MysqlContext
//	log log.Log
//}
//
//func (h *GroupAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
//		return
//	}
//	if r.Header[ADMIN_TOKEN_HEADER] == nil {
//		h.log.Info("Node add no Admin-token")
//		api.ReturnError(r, w, errors.Jerror("No Admin-token in header"), errors.BadRequestError, h.log)
//		return
//	}
//	if strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) != 0 {
//		h.log.Info("Node add token invalid")
//		api.ReturnError(r, w, errors.Jerror("Admin token unauthorized"), errors.UnauthorizedError, h.log)
//		return
//	}
//
//	result, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
//		return
//	}
//	r.Body.Close()
//
//	data := &GroupMessage{}
//	err = json.Unmarshal(result, &data)
//	if err != nil {
//		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
//		return
//	}
//	h.log.Info("Node add request: (%s) from client: %s", string(result), r.RemoteAddr)
//
//	//TODO: check args
//	if data.Ip == "" {
//		h.log.Info("Not found ip")
//		api.ReturnError(r, w, errors.Jerror("No ip in data"), errors.BadRequestError, h.log)
//		return
//	}
//
//	key := JOB_NODE_LOC + data.Ip
//	msg, err := h.eh.Get(key)
//	if err != nil {
//		h.log.Error("Get node from etcd failed")
//		api.ReturnError(r, w, errors.Jerror("check token with backend failed"), errors.BadGatewayError, h.log)
//		return
//	}
//	if msg != nil {
//		h.log.Error("Node exists, cannot add")
//		api.ReturnError(r, w, errors.Jerror("Node exists"), errors.ConflictError, h.log)
//		return
//	}
//
//	err = h.eh.Set(key, "1")
//	if err != nil {
//		h.log.Error("Set node failed")
//		api.ReturnError(r, w, errors.Jerror("Set to backend failed"), errors.BadGatewayError, h.log)
//		return
//	}
//
//	api.ReturnResponse(r, w, "", h.log)
//}
//
//func (h *GroupDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
//		return
//	}
//	if r.Header[ADMIN_TOKEN_HEADER] == nil {
//		h.log.Info("Node delete no Admin-token")
//		api.ReturnError(r, w, errors.Jerror("No Admin-token in header"), errors.BadRequestError, h.log)
//		return
//	}
//	if strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) != 0 {
//		h.log.Info("Node delete token invalid")
//		api.ReturnError(r, w, errors.Jerror("Admin token unauthorized"), errors.UnauthorizedError, h.log)
//		return
//	}
//
//	result, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
//		return
//	}
//	r.Body.Close()
//
//	data := &NodeMessage{}
//	err = json.Unmarshal(result, &data)
//	if err != nil {
//		api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
//		return
//	}
//	h.log.Info("Delete node request: (%s) from client: %s", string(result), r.RemoteAddr)
//
//	//TODO: check args
//	if data.Ip == "" {
//		h.log.Info("Not found ip")
//		api.ReturnError(r, w, errors.Jerror("No ip in data"), errors.BadRequestError, h.log)
//		return
//	}
//
//	key := JOB_NODE_LOC + data.Ip
//	msg, err := h.eh.Get(key)
//	if err != nil {
//		h.log.Error("Get node from etcd failed")
//		api.ReturnError(r, w, errors.Jerror("check token with backend failed"), errors.BadGatewayError, h.log)
//		return
//	}
//	if msg == nil {
//		h.log.Error("Node not exist")
//		api.ReturnError(r, w, errors.Jerror("Node not exist"), errors.NoContentError, h.log)
//		return
//	}
//
//	err = h.eh.UnSet(key)
//	if err != nil {
//		h.log.Error("Unset node from etcd failed")
//		api.ReturnError(r, w, errors.Jerror("Delete node from backend failed"), errors.BadGatewayError, h.log)
//		return
//	}
//
//	api.ReturnResponse(r, w, "", h.log)
//}
//
//func (h *GroupReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	var nrm NodeReadMessage
//	if r.Method != "GET" {
//		api.ReturnError(r, w, errors.Jerror("Method invalid"), errors.BadRequestError, h.log)
//		return
//	}
//	if r.Header[ADMIN_TOKEN_HEADER] == nil {
//		h.log.Info("Node read no Admin-token")
//		api.ReturnError(r, w, errors.Jerror("No Admin-token in header"), errors.BadRequestError, h.log)
//		return
//	}
//	if strings.Compare(r.Header[ADMIN_TOKEN_HEADER][0], AdminToken) != 0 {
//		h.log.Info("Node read token invalid")
//		api.ReturnError(r, w, errors.Jerror("Admin token unauthorized"), errors.UnauthorizedError, h.log)
//		return
//	}
//
//	h.log.Info("Node read request from client: %s", r.RemoteAddr)
//	//result, err := ioutil.ReadAll(r.Body)
//	//if err != nil {
//	//	api.ReturnError(r, w, errors.Jerror("Read from body failed"), errors.BadRequestError, h.log)
//	//	return
//	//}
//	//r.Body.Close()
//
//	//TODO: auth
//	//data := &NodeMessage{}
//	//err = json.Unmarshal(result, &data)
//	//if err != nil {
//	//	api.ReturnError(r, w, errors.Jerror("Parse from body failed"), errors.BadRequestError, h.log)
//	//	return
//	//}
//	//h.log.Info("Add record request: (%s) from client: %s", data, r.RemoteAddr)
//
//	////TODO: check args
//	//if data.Ip == "" {
//	//	h.log.Info("Not found ip")
//	//	api.ReturnError(r, w, errors.Jerror("No ip in data"), errors.BadRequestError, h.log)
//	//	return
//	//}
//	//if data.Group = "" {
//	//	h.log.Info("Not found group")
//	//	api.ReturnError(r, w, errors.Jerror("No group in data"), errors.BadRequestError, h.log)
//	//	return
//	//}
//
//	key := JOB_NODE_LOC
//	msg, err := h.eh.GetWithPrefix(key)
//	if err != nil {
//		h.log.Error("Get node from etcd failed")
//		api.ReturnError(r, w, errors.Jerror("check token with backend failed"), errors.BadGatewayError, h.log)
//		return
//	}
//	if msg == nil {
//		h.log.Error("Node not exist")
//		api.ReturnError(r, w, errors.Jerror("Node not exist"), errors.NoContentError, h.log)
//		return
//	}
//
//	for _, m := range msg {
//		n := &NodeMessage{}
//		n.Ip = strings.TrimPrefix(string(m.Key), JOB_NODE_LOC)
//		nrm.Nodes = append(nrm.Nodes, n)
//	}
//
//	vnrm, _ := json.Marshal(nrm)
//
//	api.ReturnResponse(r, w, string(vnrm), h.log)
//}
//
