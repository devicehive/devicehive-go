package transport

import (
	"sync"
	"fmt"
)

var mu = sync.Mutex{}

type clientsMap map[string]*client

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

	m[key] = c

	return c
}

func (m clientsMap) get(key string) (client *client, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Getting by key:", key)

	req, ok := m[key]

	if ok {
		fmt.Println("OK")
	} else {
		fmt.Println("BAD")
	}

	return req, ok
}

func (m clientsMap) forEach(f func(req *client)) {
	mu.Lock()
	defer mu.Unlock()

	for _, req := range m {
		f(req)
	}
}
