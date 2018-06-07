package transport

type Subscription struct {
	DataChan chan []byte
	ErrChan  chan error
}
