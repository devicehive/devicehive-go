package transport

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

func newWSSubscriptionsBuffer(clients *clientsMap) *wsSubscriptions {
	return &wsSubscriptions{
		clientsMap: clients,
		mu:         sync.Mutex{},
	}
}

type wsSubscriptions struct {
	*clientsMap
	buffer [][]byte
	mu     sync.Mutex
}

func (s *wsSubscriptions) BufferPut(b []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buffer = append(s.buffer, b)
}

func (s *wsSubscriptions) createSubscriber(key string) *client {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := s.clientsMap.createSubscriber(key)

	subsData, newBuffer := s.getSubscriberData(key)

	for _, b := range subsData {
		client.data <- b
	}

	s.buffer = newBuffer

	return client
}

func (s *wsSubscriptions) getSubscriberData(subsId string) (subsData [][]byte, newBuffer [][]byte) {
	ids := &ids{}
	for _, b := range s.buffer {
		err := json.Unmarshal(b, ids)
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

func (s *wsSubscriptions) CleanupBufferByTimeout(timeout time.Duration) {
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
