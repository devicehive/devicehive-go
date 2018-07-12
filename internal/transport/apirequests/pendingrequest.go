// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package apirequests

func NewPendingRequest() *PendingRequest {
	return &PendingRequest{
		Data:   make(chan []byte),
		Signal: make(chan struct{}),
		Err:    make(chan error),
	}
}

type PendingRequest struct {
	Data   chan []byte
	Err    chan error
	Signal chan struct{}
}

func (r *PendingRequest) Close() {
	close(r.Signal)
}
