package transport

import (
	"sync"
)

var mu = sync.Mutex{}

type requestMap map[string]*request

func (m requestMap) delete(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(m, key)
}

func (m requestMap) create(key string) (req *request) {
	mu.Lock()
	defer mu.Unlock()

	res, err := make(chan []byte), make(chan *Error)

	m[key] = &request{
		response: res,
		err:      err,
	}

	return m[key]
}

func (m requestMap) get(key string) (request *request, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	req, ok := m[key]
	return req, ok
}

func (m requestMap) forEach(f func(req *request)) {
	mu.Lock()
	defer mu.Unlock()

	for _, req := range m {
		f(req)
	}
}
