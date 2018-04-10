package transport

import (
	"sync"
)

var mu = sync.Mutex{}

type clientsMap map[string]*client

type client struct {
	data chan []byte
	err  chan *Error
}

func (c *client) close() {
	close(c.data)
	close(c.err)
}

func (m clientsMap) delete(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(m, key)
}

func (m clientsMap) createClient(key string) (req *client) {
	return m.create(key, true)
}

func (m clientsMap) createSubscriber(key string) (req *client) {
	return m.create(key, false)
}

func (m clientsMap) create(key string, isErrChan bool) (req *client) {
	mu.Lock()
	defer mu.Unlock()

	var c *client
	res := make(chan []byte)
	if isErrChan {
		err := make(chan *Error)
		c = &client{
			data: res,
			err:  err,
		}
	} else {
		c = &client{
			data: res,
		}
	}

	m[key] = c

	return c
}

func (m clientsMap) get(key string) (client *client, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	req, ok := m[key]
	return req, ok
}

func (m clientsMap) forEach(f func(req *client)) {
	mu.Lock()
	defer mu.Unlock()

	for _, req := range m {
		f(req)
	}
}
