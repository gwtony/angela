package handler

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

type ReportMessage struct {
	Meta   OrchMeta       `json:"meta"`
	State  []ReportState  `json:"state"`
	Output []ReportOutput `json:"output"`
}

type ReportState struct {
	Node   string `json:"node"`
	State string  `json:"state"`
	Rc     int    `json:"rc"`
}
type ReportOutput struct {
	Node   string `json:"node"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Start  string `json:"start"`
	Stop   string `json:"stop"`
	Delta  uint64 `json:"delta"`
}

type Node struct {
	Node   string `json:"node"`
	//State string `json:"state"`
}

