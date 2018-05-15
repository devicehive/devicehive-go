package transport

import (
	"sync"
)

func newClientsMap() *clientsMap {
	return &clientsMap{
		clients: make(map[string]*client),
		mu: sync.Mutex{},
	}
}

type clientsMap struct {
	clients map[string]*client
	mu sync.Mutex
}

type client struct {
	data   chan []byte
	err    chan *Error
	signal chan struct{}
}

func (c *client) close() {
	close(c.data)
	close(c.signal)

	if c.err != nil {
		close(c.err)
	}
}

func (m *clientsMap) delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, key)
}

func (m *clientsMap) createClient(key string) (req *client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, true)
}

func (m *clientsMap) createSubscriber(key string) (req *client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, false)
}

func (m *clientsMap) create(key string, isErrChan bool) (req *client) {
	var c *client
	res := make(chan []byte, 16)
	signal := make(chan struct{})
	if isErrChan {
		err := make(chan *Error)
		c = &client{
			data:   res,
			err:    err,
			signal: signal,
		}
	} else {
		c = &client{
			data:   res,
			signal: signal,
		}
	}

	m.clients[key] = c

	return c
}

func (m *clientsMap) get(key string) (client *client, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	req, ok := m.clients[key]
	return req, ok
}

func (m *clientsMap) forEach(f func(req *client)) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, req := range m.clients {
		f(req)
	}
}
