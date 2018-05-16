package apirequests

type PendingRequest struct {
	Data   chan []byte
	Err    chan error
	Signal chan struct{}
}

func (c *PendingRequest) Close() {
	close(c.Data)
	close(c.Signal)

	if c.Err != nil {
		close(c.Err)
	}
}
