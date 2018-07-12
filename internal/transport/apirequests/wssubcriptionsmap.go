// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

import (
	"strconv"
	"sync"

	"github.com/devicehive/devicehive-go/internal/utils"
)

func NewWSSubscriptionsMap() *WSSubscriptionsMap {
	return &WSSubscriptionsMap{
		subscriptions: make(map[string]*WSSubscription),
	}
}

type WSSubscriptionsMap struct {
	subscriptions       map[string]*WSSubscription
	subscriptionsLocker sync.RWMutex
	buffer              [][]byte
	bufferLocker        sync.RWMutex
	mu                  sync.RWMutex
}

func (s *WSSubscriptionsMap) Get(key string) (*WSSubscription, bool) {
	s.subscriptionsLocker.RLock()
	wssub, ok := s.subscriptions[key]
	s.subscriptionsLocker.RUnlock()

	return wssub, ok
}

func (s *WSSubscriptionsMap) Delete(key string) {
	s.subscriptionsLocker.Lock()
	defer s.subscriptionsLocker.Unlock()

	delete(s.subscriptions, key)
}

func (s *WSSubscriptionsMap) ForEach(f func(req *WSSubscription)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.subscriptions {
		f(sub)
	}
}

func (s *WSSubscriptionsMap) BufferPut(b []byte) {
	s.bufferLocker.Lock()
	s.buffer = append(s.buffer, b)
	s.bufferLocker.Unlock()
}

func (s *WSSubscriptionsMap) CreateSubscription(key string) *WSSubscription {
	subs := &WSSubscription{
		PendingRequest: NewPendingRequest(),
	}
	s.subscriptionsLocker.Lock()
	s.subscriptions[key] = subs
	s.subscriptionsLocker.Unlock()

	subsData, newBuffer := s.extractSubscriberData(key)

	go func() {
		for _, b := range subsData {
			subs.Data <- b
		}
	}()

	s.bufferLocker.Lock()
	s.buffer = newBuffer
	s.bufferLocker.Unlock()

	return subs
}

func (s *WSSubscriptionsMap) extractSubscriberData(subsId string) (subsData [][]byte, newBuffer [][]byte) {
	s.bufferLocker.RLock()
	defer s.bufferLocker.RUnlock()
	for _, b := range s.buffer {
		ids, err := utils.ParseIDs(b)
		if err != nil {
			continue
		}

		dataSubsId := strconv.FormatInt(ids.Subscription, 10)

		if dataSubsId != "" && dataSubsId == subsId {
			subsData = append(subsData, b)
		} else if dataSubsId != "" && dataSubsId != subsId {
			newBuffer = append(newBuffer, b)
		}
	}

	return
}
