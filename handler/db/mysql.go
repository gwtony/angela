package db

//import (
//	"fmt"
//	//"time"
//	"database/sql"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/gwtony/gapi/log"
//	"github.com/gwtony/gapi/errors"
//	"github.com/gwtony/angela/handler/variable"
//)
//
//type MysqlContext struct {
//	addr    string
//	dbname  string
//	dbuser  string
//	dbpwd   string
//
//	db      *sql.DB
//	login   string
//
//	log     log.Log
//}
//
//const (
//	//DCRON_INSERT_SQL = "insert into results (run_id, result, create_time) values (?, ?, ?)"
//	//DCRON_READ_SQL   = "select result from results where run_id = ?"
//	//DCRON_UPDATE_SQL = "update results set result = ?, create_time = ? where run_id = ?"
//)
//
//var mci *MysqlContext
//
//func InitMysqlContext(addr, dbname, dbuser, dbpwd string, log log.Log) error {
//	mc := &MysqlContext{}
//
//	mc.log      = log
//	mc.addr     = addr
//	mc.dbname   = dbname
//	mc.dbuser   = dbuser
//	mc.dbpwd    = dbpwd
//	mc.login    = fmt.Sprintf("%s:%s@tcp(%s)/%s", dbuser, dbpwd, addr, dbname)
//	db, err     := sql.Open("mysql", mc.login)
//	if err != nil {
//		mc.log.Error("Open mysql context failed: %s", err)
//		return err
//	}
//	mc.db = db
//	mci = mc
//
//	return nil
//}
//
//func (mc *MysqlContext) Close() error {
//	return mc.db.Close()
//}
//
//func (mc *MysqlContext) QueryRead(runid string) (string, error) {
//	var result string
//
//	rows, err := mc.db.Query(variable.DCRON_READ_SQL, runid)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			mc.log.Error("Scan no answer")
//			return "", errors.NoContentError
//		}
//		mc.log.Error("Execute get result for runid %s failed: %s", runid, err)
//		return "", errors.BadGatewayError
//	}
//	defer rows.Close()
//
//	result = ""
//
//	for rows.Next() {
//		err := rows.Scan(&result)
//		if err == sql.ErrNoRows {
//			mc.log.Error("Scan no answer")
//			return "", errors.NoContentError
//		}
//		if err != nil {
//			mc.log.Error("Scan read answer failed: %s", err)
//			return "", errors.InternalServerError
//		}
//		//Should only one record
//		break
//	}
//
//	return result, nil
//}
//
//func (mc *MysqlContext) QueryInsert(runid, result string, create_time int) error {
//	return mc.QueryWrite(variable.DCRON_INSERT_SQL, runid, result, create_time)
//}
//
//func (mc *MysqlContext) QueryWrite(query string, args ...interface{}) error {
//	res, err := mc.db.Exec(query, args...)
//
//	if err != nil {
//		mc.log.Error("Execute write sql: ", query, args, " failed: ", err)
//		return errors.BadGatewayError
//	}
//	affected, err := res.RowsAffected()
//	if err != nil {
//		mc.log.Error("Get rows affected failed: %s", err)
//		return errors.InternalServerError
//	}
//	if int(affected) <= 0 {
//		return errors.BadGatewayError
//	}
//
//	return nil
//}
//
//func (mc *MysqlContext) QueryUpdate(runid, result string, create_time int) error {
//	if runid == "" {
//		return errors.BadRequestError
//	}
//
//	res, err := mc.db.Exec(variable.DCRON_UPDATE_SQL, result, create_time, runid)
//
//	if err != nil {
//		mc.log.Error("Execute update result for runid: %s failed: %s", runid, err)
//		return errors.BadGatewayError
//	}
//
//	affected, err := res.RowsAffected()
//	if err != nil {
//		mc.log.Error("Get rows affected failed: %s", err)
//		return errors.InternalServerError
//	}
//
//	if int(affected) > 0 {
//		return nil
//	}
//
//	mc.log.Info("No such runid: %s", runid)
//	return errors.NoContentError
//}
//
//func Close() error {
//	return mci.Close()
//}
//
//func QueryUpdate(runid, result string, create_time int) error {
//	return mci.QueryUpdate(runid, result, create_time)
//}
//
//func QueryRead(runid string) (string, error) {
//}
//func QueryInsert(runid, result string, create_time int) error {
//	return mc.QueryWrite(variable.DCRON_INSERT_SQL, runid, result, create_time)
//}
//
//func QueryWrite(query string, args ...interface{}) error {
//}
