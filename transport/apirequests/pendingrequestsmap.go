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

func (m *PendingRequestsMap) CreateRequest(key string) (req *PendingRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, true)
}

func (m *PendingRequestsMap) CreateSubscription(key string) (req *PendingRequest) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.create(key, true)
}

func (m *PendingRequestsMap) create(key string, isErrChan bool) (req *PendingRequest) {
	var c *PendingRequest
	data := make(chan []byte, 16)
	signal := make(chan struct{})

	c = &PendingRequest{
		Data:       data,
		Signal:     signal,
		DataLocker: sync.Mutex{},
	}
	if isErrChan {
		c.Err = make(chan error)
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
