package transport

import (
	"sync"
)

var mu = sync.Mutex{}

type clientsMap map[string]*client

type client struct {
	response chan []byte
	err      chan *Error
}

func (c *client) close() {
	close(c.response)
	close(c.err)
}

func (m clientsMap) delete(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(m, key)
}

func (m clientsMap) create(key string) (req *client) {
	mu.Lock()
	defer mu.Unlock()

	res, err := make(chan []byte), make(chan *Error)

	m[key] = &client{
		response: res,
		err:      err,
	}

	return m[key]
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
