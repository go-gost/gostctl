package api

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
