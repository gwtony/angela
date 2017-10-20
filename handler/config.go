package handler

import (
	"fmt"
	"os"
	"time"
	"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
)

type AngelaConfig struct {
	adminToken    string        /* admin token */

	eaddr         []string      /* etcd addr */
	eto           time.Duration /* etcd timeout */

	apiLoc        string        /* angela api location */

	orchTo        time.Duration
	orchUrl       string        /* orch url */
	orchCreateLoc string        /* orch create location */
	orchCancelLoc string        /* orch cancel location */
	orchStateLoc  string        /* orch state location */
	orchTeeLoc    string        /* orch tee location */
	orchGroup     string        /* orch group */
	orchToken     string        /* orch token */
	orchCb        string        /* orch call back */

	raddr         string        /* redis addr to send alert */
	alertEnable   bool          /* alert enable */

	maddr         string        /* mysql addr */
	dbname        string        /* db name */
	dbuser        string        /* db username */
	dbpwd         string        /* db password */

	name          string        /* angela server name */

	timeout       time.Duration
}

// ParseConfig parses config
func (conf *AngelaConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}
	conf.maddr, err = cf.C.GetString("angela", "mysql_addr")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_addr")
		return err
	}
	conf.dbname, err = cf.C.GetString("angela", "mysql_dbname")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbname")
		return err
	}
	conf.dbuser, err = cf.C.GetString("angela", "mysql_dbuser")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbuser")
		return err
	}
	conf.dbpwd, err = cf.C.GetString("angela", "mysql_dbpwd")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbpwd")
		return err
	}

	conf.adminToken, err = cf.C.GetString("angela", "admin_token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [angela] Read conf: No admin_token, use default admin token:", DEFAULT_ADMIN_TOKEN)
		conf.adminToken = DEFAULT_ADMIN_TOKEN
	}

	return nil
}
