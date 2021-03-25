package p2p_matrix

type ScriptModel struct {
	Name        string `json:"name"`
	Operations  int    `json:"operations"`
	NodesAmount int    `json:"nodesAmount"`
	PeersMin    int    `json:"peersMin"`
	PeersMax    int    `json:"peersMax"`
	FileSizeMin int    `json:"fileSizeMin"`
	FileSizeMax int    `json:"fileSizeMax"`

	Nodes map[int]Node   `json:"nodes"`
	Files map[int]File   `json:"files"`
	Story []StoryElement `json:"story"`
	Pings Pings          `json:"pings"`
}

type Node struct {
	Location int     `json:"location"`
	Speed    float64 `json:"speed"`
}

type File struct {
	Size float64 `json:"size"`
}

type StoryElement struct {
	NodeActions []StoryElementAction    `json:"nodeActions"`
	Operations  []StoryElementOperation `json:"operations"`
}

type StoryElementOperation struct {
	Node int    `json:"nodeId"`
	File int    `json:"fileId"`
	Type string `json:"type"`
}

type StoryElementAction struct {
	Node   int    `json:"nodeId"`
	Action string `json:"action"`
}

type Pings struct {
	Locations map[int]string `json:"locations"`
	Pings     []pingPair     `json:"pings"`
}

// todo optimization checkbox
// todo why float32 not float64
type pingPair struct {
	Location1 int     `json:"l1"`
	Location2 int     `json:"l2"`
	Ping      float64 `json:"p"`
}
