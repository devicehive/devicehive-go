// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

import "sync"

type PendingRequest struct {
	Data       chan []byte
	Err        chan error
	Signal     chan struct{}
	DataLocker sync.Mutex
}

func (r *PendingRequest) Close() {
	r.DataLocker.Lock()
	defer r.DataLocker.Unlock()
	close(r.Signal)
}
