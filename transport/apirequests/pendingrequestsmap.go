// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

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

func (m *PendingRequestsMap) CreateRequest(key string) *PendingRequest {
	req := &PendingRequest{
		Data:       make(chan []byte, 16),
		Signal:     make(chan struct{}),
		DataLocker: sync.Mutex{},
		Err:        make(chan error),
	}

	m.mu.Lock()
	m.clients[key] = req
	m.mu.Unlock()

	return req
}

func (m *PendingRequestsMap) Get(key string) (*PendingRequest, bool) {
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
