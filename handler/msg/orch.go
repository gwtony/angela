package msg

type OrchMeta struct {
	Runid string `json:"runid"`
}
type OrchMessage struct {
	Meta     OrchMeta `json:"meta"`
	Nodes    []string `json:"nodes"`
	Command  string   `json:"command"`
	User     string   `json:"user"`
	Timeout  uint64   `json:"timeout"`
	Parallel int      `json:"parallel"`
	Callback string   `json:"callback"`
	Type     string   `json:"type"`
}

type OrchMapMessage struct {
	Jobid     string `json:"jobid"` //cron jobid
	Timestamp uint64 `json:"timestamp"`
	Expire    uint64 `json:"expire"`
	Once      int    `json:"once"`
	Parent    string `json:"parent"`
	ParentTs  uint64 `json:"parentts"`
	Event     string `json:"event"`
}

