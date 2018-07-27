package authmanager

import (
	"sync"
	"time"
)

type ReauthenticationState struct {
	lastReauth      time.Time
	lastReauthMutex sync.RWMutex
}

func (sr *ReauthenticationState) Needed() bool {
	sr.lastReauthMutex.RLock()
	defer sr.lastReauthMutex.RUnlock()
	return time.Now().Sub(sr.lastReauth) > 5*time.Second
}

func (sr *ReauthenticationState) Checkpoint() {
	sr.lastReauthMutex.Lock()
	defer sr.lastReauthMutex.Unlock()
	sr.lastReauth = time.Now()
}
