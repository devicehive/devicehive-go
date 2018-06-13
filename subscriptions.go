package devicehive_go

import (
	"sync"
	"time"
)

var subscriptionReauth = &subscriptionReauthenticator{}

type subscriber interface {
	Remove() *Error
	sendError(*Error)
}

func removeSubscriptionWithError(s subscriber, err *Error) {
	s.sendError(err)
	rmErr := s.Remove()
	if rmErr != nil {
		s.sendError(rmErr)
	}
}

type subscriptionReauthenticator struct {
	lastReauth      time.Time
	lastReauthMutex sync.RWMutex
}

func (sr *subscriptionReauthenticator) reauthNeeded() bool {
	sr.lastReauthMutex.RLock()
	defer sr.lastReauthMutex.RUnlock()
	return time.Now().Sub(sr.lastReauth) > 5*time.Second
}

func (sr *subscriptionReauthenticator) reauthPoint() {
	sr.lastReauthMutex.Lock()
	defer sr.lastReauthMutex.Unlock()
	sr.lastReauth = time.Now()
}
