package handler

import (
	"sync"
	"github.com/gwtony/gapi/log"
	"github.com/gwtony/gapi/api"
	//"github.com/gwtony/gapi/hserver"
	"github.com/gwtony/gapi/config"
	//"github.com/gwtony/gapi/errors"
)

//var DMapLock sync.Mutex
//var DcronMap map[string]*CronMessage

var AdminToken string

// InitContext inits angela context
func InitContext(conf *config.Config, log log.Log) error {
	var rh *RedisHandler
	cf := &AngelaConfig{}
	err := cf.ParseConfig(conf)
	if err != nil {
		log.Error("Angela parse config failed")
		return err
	}
	AdminToken = cf.adminToken

	mc, err := InitMysqlContext(cf.maddr, cf.dbname, cf.dbuser, cf.dbpwd, log)
	if err != nil {
		log.Error("Dcron init mysql context failed")
		return err
	}

	//DcronMap = make(map[string]*CronMessage, 10000)

	//eh := InitEtcdHandler(cf.eaddr, cf.eto, log)

	//cw := &CronWorker{eh: eh, name: cf.name, log: log}

	//ch := NewCronHandler(log)

	//watcher := InitWatcher(ch, cw, eh, log)
	//elector := &Elector{name: cf.name, eh: eh, interval: 3, term: 6, log: log}

	//if err := elector.Register(); err != nil {
	//	log.Error("Another dcron named: %s is running", cf.name)
	//	return err
	//}
	//go elector.Run()

	//orch_url := cf.orchUrl
	//crloc    := cf.orchCreateLoc
	//caloc    := cf.orchCancelLoc
	//sloc     := cf.orchStateLoc
	//tloc     := cf.orchTeeLoc
	//group    := cf.orchGroup
	//token    := cf.orchToken

	//cb       := "http://" + cf.orchCb + DCRON_LOC + DCRON_REPORT_LOC

	//InitOrchHandler(orch_url, crloc, caloc, sloc, tloc, group, token, cb ,log)

	//apiLoc := cf.apiLoc

	////log.Debug(apiLoc, DCRON_ADD_LOC)
	//if cf.alertEnable {
	//	rh = InitRedisHandler(cf.raddr, log)
	//	go rh.Run()
	//} else {
	//	rh = nil
	//}

	////collect jobs first
	//err = watcher.PullAll(JOB_META_LOC)
	//if err != nil {
	//	log.Error("Pull jobs from etcd failed:", err)
	//	return err
	//}

	//go watcher.Run()
	//go MonitorRun(rh, log)

	api.AddHttpHandler(apiLoc + JOB_CREATE_LOC, &JobCreateHandler{mc: mc, log: log})
	api.AddHttpHandler(apiLoc + JOB_CANCEL_LOC, &JobCancelHandler{mc: mc, log: log})

	api.AddHttpHandler(apiLoc + GROUP_ADD_LOC, &GroupAddHandler{mc: mc, log: log})
	api.AddHttpHandler(apiLoc + GROUP_DELETE_LOC, &GroupDeleteHandler{mc: mc, log: log})
	api.AddHttpHandler(apiLoc + GROUP_READ_LOC, &GroupReadHandler{mc: mc, log: log})

	ch.Run()

	return nil
}
