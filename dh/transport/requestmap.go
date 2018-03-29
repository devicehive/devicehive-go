package transport

import "sync"

var mu = sync.Mutex{}

type requestMap map[string]chan devicehiveData

func (m requestMap) delete(key string) {
	mu.Lock()
	defer mu.Unlock()

	delete(m, key)
}

func (m requestMap) add(key string, val chan devicehiveData) {
	mu.Lock()
	defer mu.Unlock()

	m[key] = val
}

func (m requestMap) get(key string) (responseChan chan devicehiveData, ok bool) {
	mu.Lock()
	defer mu.Unlock()

	resChan, ok := m[key]
	return resChan, ok
}