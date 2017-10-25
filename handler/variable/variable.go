package variable

const (
	// VERSION version
	VERSION                 = "0.1 alpha"

	API_CONTENT_HEADER      = "application/json;charset=utf-8"

	DEFAULT_ADMIN_TOKEN     = "ORCH_TOKEN"
	ADMIN_TOKEN_HEADER      = "Admin-Token"

	//ETCD_CONTENT_HEADER      = "application/x-www-form-urlencoded"
	ORCH_LOC               = "/orch"

	//http location for job
	JOB_CREATE_LOC          = "/job/create"
	JOB_CANCEL_LOC          = "/job/cancel"

	//http location for group
	GROUP_ADD_LOC           = "/group/add"
	GROUP_DELETE_LOC        = "/group/delete"
	GROUP_READ_LOC          = "/group/read"

	//orch job status
	JOB_LAUNCHED            = "launched"
	JOB_RUNNING             = "running"
	JOB_SUCCESS             = "success"
	JOB_FAILED              = "failed"
	JOB_ABORT               = "abort"
	JOB_TIMEOUT             = "timeout"

	DCRON_TOKEN             = "dcron_token"

	//operation in orch handler
	ORCH_CREATE             = iota
	ORCH_CANCEL
	ORCH_READ
	ORCH_GROUP_HEADER       = "ORCH-AUTH-GROUP"
	ORCH_TOKEN_HEADER       = "ORCH-AUTH-TOKEN"
	ORCH_RETRY_NUM          = 3
	ORCH_OPERATE_TYPE       = "exec"

	SUBJOB_PADDING          = "padding"
	ORCH_RUNTIME_DIR        = "/tmp/orch_run"
	ORCH_TMP_DIR            = "/tmp/orch_tmp"
)
