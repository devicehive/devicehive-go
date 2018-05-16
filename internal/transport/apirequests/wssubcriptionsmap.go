package apirequests

import (
	"strconv"
	"sync"
	"time"
	"github.com/devicehive/devicehive-go/internal/utils"
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

	client := s.PendingRequestsMap.CreateSubscription(key)

	subsData, newBuffer := s.getSubscriberData(key)

	for _, b := range subsData {
		client.Data <- b
	}

	s.buffer = newBuffer

	return client
}

func (s *WSSubscriptionsMap) getSubscriberData(subsId string) (subsData [][]byte, newBuffer [][]byte) {
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

	return subsData, newBuffer
}

func (s *WSSubscriptionsMap) CleanupBufferByTimeout(timeout time.Duration) {
	for {
		time.Sleep(timeout)

		s.mu.Lock()
		l := len(s.buffer)
		if l == 0 {
			s.mu.Unlock()
			continue
		}

		m := (l - 1) / 2
		s.buffer = s.buffer[m:]
		s.mu.Unlock()
	}
}
