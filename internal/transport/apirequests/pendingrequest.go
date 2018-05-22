package apirequests

import "sync"

type PendingRequest struct {
	Data   chan []byte
	Err    chan error
	Signal chan struct{}
	DataLocker sync.Mutex
}

func (c *PendingRequest) Close() {
	c.DataLocker.Lock()
	defer c.DataLocker.Unlock()
	close(c.Data)
	close(c.Signal)

	if c.Err != nil {
		close(c.Err)
	}
}
