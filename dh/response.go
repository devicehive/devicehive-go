package dh

type response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

type httpResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
