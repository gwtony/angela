package msg

type NodeState struct {
	Node    string `json:"node"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Rc      int    `json:rc`
	Timeout bool   `json:"timeout"`
	Error   error  `json:"error"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Delta   uint64 `json:"delta"`
}
