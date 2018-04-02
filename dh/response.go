package dh

type response struct {
	Action string `json:"action"`
	Status string `json:"status"`
	Error string `json:"error"`
	Code int `json:"code"`
}
