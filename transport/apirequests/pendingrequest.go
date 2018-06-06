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

func (c *PendingRequest) Close() {
	c.DataLocker.Lock()
	defer c.DataLocker.Unlock()
	close(c.Data)
	close(c.Signal)

	if c.Err != nil {
		close(c.Err)
	}
}
