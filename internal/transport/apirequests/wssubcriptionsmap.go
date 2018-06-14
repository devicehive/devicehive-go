// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

import (
	"strconv"
	"sync"

	"github.com/devicehive/devicehive-go/internal/utils"
)

func NewWSSubscriptionsMap(clients *PendingRequestsMap) *WSSubscriptionsMap {
	return &WSSubscriptionsMap{
		PendingRequestsMap: clients,
	}
}

type WSSubscriptionsMap struct {
	*PendingRequestsMap
	buffer [][]byte
	mu     sync.RWMutex
}

func (s *WSSubscriptionsMap) BufferPut(b []byte) {
	s.mu.Lock()
	s.buffer = append(s.buffer, b)
	s.mu.Unlock()
}

func (s *WSSubscriptionsMap) CreateSubscription(key string) *PendingRequest {
	subs := s.PendingRequestsMap.CreateRequest(key)

	subsData, newBuffer := s.extractSubscriberData(key)

	go func() {
		for _, b := range subsData {
			subs.Data <- b
		}
	}()

	s.mu.Lock()
	s.buffer = newBuffer
	s.mu.Unlock()

	return subs
}

func (s *WSSubscriptionsMap) extractSubscriberData(subsId string) (subsData [][]byte, newBuffer [][]byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
