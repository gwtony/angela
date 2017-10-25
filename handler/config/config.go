package config

import (
	"fmt"
	"os"
	"time"
	//"strings"
	"github.com/gwtony/gapi/config"
	"github.com/gwtony/gapi/errors"
	"github.com/gwtony/angela/handler/variable"
)

type AngelaConfig struct {
	AdminToken    string        /* admin token */

	ApiLoc        string        /* angela api location */

	Maddr         string        /* mysql addr */
	Dbname        string        /* db name */
	Dbuser        string        /* db username */
	Dbpwd         string        /* db password */

	SshKey        string        /* ssh key */

	Timeout       time.Duration
}

// ParseConfig parses config
func (conf *AngelaConfig) ParseConfig(cf *config.Config) error {
	var err error
	if cf.C == nil {
		return errors.BadConfigError
	}

	//conf.Maddr, err = cf.C.GetString("angela", "mysql_addr")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_addr")
	//	return err
	//}
	//conf.Dbname, err = cf.C.GetString("angela", "mysql_dbname")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbname")
	//	return err
	//}
	//conf.Dbuser, err = cf.C.GetString("angela", "mysql_dbuser")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbuser")
	//	return err
	//}
	//conf.Dbpwd, err = cf.C.GetString("angela", "mysql_dbpwd")
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No mysql_dbpwd")
	//	return err
	//}

	conf.AdminToken, err = cf.C.GetString("angela", "admin_token")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [angela] Read conf: No admin_token, use default admin token:", variable.DEFAULT_ADMIN_TOKEN)
		conf.AdminToken = variable.DEFAULT_ADMIN_TOKEN
	}

	conf.ApiLoc, err = cf.C.GetString("angela", "api_loc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [angela] Read conf: No api_loc, use default loc:", variable.ORCH_LOC)
		conf.ApiLoc = variable.ORCH_LOC
	}

	conf.SshKey, err = cf.C.GetString("angela", "ssh_key")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] [angela] Read conf: No ssh_key")
		return err
	}

	return nil
}
