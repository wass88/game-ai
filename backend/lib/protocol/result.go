package protocol

type ResultA struct {
	Record    string `json:"record"`
	Exception string `json:"exception"`
}

type ResultPlayerA struct {
	Result    int    `json:"result"`
	Stderr    string `json:"stderr"`
	Exception string `json:"exception"`
}

type Result struct {
	Game      string
	Record    []string
	Exception string
	Result    []ResultPlayer
}

type ResultPlayer struct {
	Result    int
	Stderr    string
	Exception string
}
