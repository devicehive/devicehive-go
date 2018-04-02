package transport

type response struct {
	data chan []byte
	err  chan *Error
}

func (r *response) close() {
	close(r.data)
	close(r.err)
}
