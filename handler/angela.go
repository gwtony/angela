package handler

import (
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	"github.com/gwtony/gapi/config"

	//"github.com/gwtony/angela/handler/db"
	aconf "github.com/gwtony/angela/handler/config"
	hh "github.com/gwtony/angela/handler/http_handler"
	"github.com/gwtony/angela/handler/variable"
	"github.com/gwtony/angela/handler/worker"
)

var AdminToken string

// InitContext inits angela context
func InitContext(conf *config.Config, log log.Log) error {
	cf := &aconf.AngelaConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Angela parse config failed")
		return err
	}
	AdminToken = cf.AdminToken

	//err = db.InitMysqlContext(cf.Maddr, cf.Dbname, cf.Dbuser, cf.Dbpwd, log)
	//if err != nil {
	//	log.Error("Angela init mysql context failed")
	//	return err
	//}
	err = worker.InitWorker(cf.SshKey)
	if err != nil {
		log.Error("Init worker failed")
		return err
	}

	api.AddHttpHandler(cf.ApiLoc + variable.JOB_CREATE_LOC, hh.InitJobCreateHandler(log))
	api.AddHttpHandler(cf.ApiLoc + variable.JOB_CANCEL_LOC, hh.InitJobCancelHandler(log))

	//api.AddHttpHandler(apiLoc + variable.GROUP_ADD_LOC, &hh.GroupAddHandler{log: log})
	//api.AddHttpHandler(apiLoc + variable.GROUP_DELETE_LOC, &hh.GroupDeleteHandler{log: log})
	//api.AddHttpHandler(apiLoc + variable.GROUP_READ_LOC, &hh.GroupReadHandler{log: log})

	return nil
}
