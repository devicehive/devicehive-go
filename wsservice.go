// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
)

// Creates low-level WS API which sends requests concurrently and writes all responses to a single channel.
// This might be useful in case of non-blocking writes (i.e. sending sensor data, subscribing for commands).
func WSConnect(url string) (*WSClient, *Error) {
	tsp, tspErr := transport.Create(url)
	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	if tsp.IsHTTP() {
		return nil, &Error{name: WrongURLErr, reason: "ws:// protocol is required"}
	}

	c := &WSClient{
		transportAdapter: transportadapter.New(tsp).(*transportadapter.WSAdapter),
		DataChan:         make(chan []byte),
		ErrorChan:        make(chan error),
	}

	return c, nil
}
