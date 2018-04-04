package transport

import (
	"sync"
)

var mu = sync.Mutex{}

type requestMap map[string]*response

func (m requestMap) delete(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(m, key)
}

func (m requestMap) create(key string) (dataChan chan []byte, errChan chan *Error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := make(chan []byte), make(chan *Error)

	m[key] = &response{
		data: data,
		err:  err,
	}

	return data, err
}

func (m requestMap) get(key string) (response *response, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	res, ok := m[key]
	return res, ok
}

func (m requestMap) forEach(f func(res *response)) {
	mu.Lock()
	defer mu.Unlock()

	for _, res := range m {
		f(res)
	}
}
