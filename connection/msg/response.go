package msg

const (
	LogResponseMsg   string = "log"
	WatchResponseMsg        = "watch"
	ErrorResponseMsg        = "error"
)

type ResponseError struct {
	Msg string `json:"message"`
}

type ErrorContent struct {
	OriginalMsg interface{}   `json:"originalMessage"`
	Error       ResponseError `json:"error"`
}

type WatchContent struct {
	RunAll bool     `json:"runAll"`
	Run    []string `json:"run"`
}

type LogContent struct {
	Msg string `json:"msg"`
}

type ResponseMsgType struct {
	Type string `json:"type"`
}
