package transport

type Subscription struct {
	DataChan chan []byte
	ErrChan  chan error
	signal   chan struct{}
}

func (s *Subscription) ContinuePolling() {
	s.signal <- struct{}{}
}
