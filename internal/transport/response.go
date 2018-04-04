package transport

type request struct {
	response chan []byte
	err      chan *Error
}

func (r *request) close() {
	close(r.response)
	close(r.err)
}
