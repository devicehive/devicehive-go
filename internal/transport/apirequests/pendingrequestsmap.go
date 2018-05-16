package apirequests

import (
	"sync"
)

func NewClientsMap() *PendingRequestsMap {
	return &PendingRequestsMap{
		clients: make(map[string]*PendingRequest),
		mu:      sync.Mutex{},
	}
}

type PendingRequestsMap struct {
	clients map[string]*PendingRequest
	mu      sync.Mutex
}

func (m *PendingRequestsMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, key)
}

func (m *PendingRequestsMap) CreateClient(key string) (req *PendingRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, true)
}

func (m *PendingRequestsMap) CreateSubscriber(key string) (req *PendingRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, false)
}

func (m *PendingRequestsMap) create(key string, isErrChan bool) (req *PendingRequest) {
	var c *PendingRequest
	res := make(chan []byte, 16)
	signal := make(chan struct{})
	if isErrChan {
		err := make(chan error)
		c = &PendingRequest{
			Data:   res,
			Err:    err,
			Signal: signal,
		}
	} else {
		c = &PendingRequest{
			Data:   res,
			Signal: signal,
		}
	}

	m.clients[key] = c

	return c
}

func (m *PendingRequestsMap) Get(key string) (client *PendingRequest, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	req, ok := m.clients[key]
	return req, ok
}

func (m *PendingRequestsMap) ForEach(f func(req *PendingRequest)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, req := range m.clients {
		f(req)
	}
}
