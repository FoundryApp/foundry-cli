package msg

type WatchfnContent struct {
	RunAll bool     `json:"runAll"`
	Run    []string `json:"run"`
}

type WatchfnMsg struct {
	Type    string         `json:"type"`
	Content WatchfnContent `json:"content"`
}

func NewWatchfnMsg(all bool, fns []string) *WatchfnMsg {
	c := WatchfnContent{all, fns}
	return &WatchfnMsg{"watch", c}
}

func (wm *WatchfnMsg) Body() interface{} {
	return wm
}
