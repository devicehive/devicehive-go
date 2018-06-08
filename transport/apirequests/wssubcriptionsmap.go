// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

import (
	"strconv"
	"sync"

	"github.com/devicehive/devicehive-go/utils"
)

func NewWSSubscriptionsMap(clients *PendingRequestsMap) *WSSubscriptionsMap {
	return &WSSubscriptionsMap{
		PendingRequestsMap: clients,
		mu:                 sync.Mutex{},
	}
}

type WSSubscriptionsMap struct {
	*PendingRequestsMap
	buffer [][]byte
	mu     sync.Mutex
}

func (s *WSSubscriptionsMap) BufferPut(b []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffer = append(s.buffer, b)
}

func (s *WSSubscriptionsMap) CreateSubscription(key string) *PendingRequest {
	s.mu.Lock()
	defer s.mu.Unlock()

	subs := s.PendingRequestsMap.CreateRequest(key)

	subsData, newBuffer := s.extractSubscriberData(key)

	go func() {
		subs.DataLocker.Lock()
		defer subs.DataLocker.Unlock()
		for _, b := range subsData {
			subs.Data <- b
		}
	}()

	s.buffer = newBuffer

	return subs
}

func (s *WSSubscriptionsMap) extractSubscriberData(subsId string) (subsData [][]byte, newBuffer [][]byte) {
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
