package msg

type ReportMessage struct {
	Meta   OrchMeta       `json:"meta"`
	State  []ReportState  `json:"state"`
	Output []ReportOutput `json:"output"`
}

type ReportState struct {
	Node   string `json:"node"`
	State  string `json:"state"`
	Rc     int    `json:"rc"`
	Error  string `json:"error"`
}
type ReportOutput struct {
	Node   string `json:"node"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Start  string `json:"start"`
	Stop   string `json:"stop"`
	Delta  uint64 `json:"delta"`
}

//type Node struct {
//	Node   string `json:"node"`
//	//State string `json:"state"`
//}

