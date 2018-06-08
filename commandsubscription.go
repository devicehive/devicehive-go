// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"sync"
	"github.com/devicehive/devicehive-go/transport"
)

var commandSubsMutex = sync.Mutex{}
var commandSubscriptions = make(map[*CommandSubscription]string)

type CommandSubscription struct {
	CommandsChan chan *Command
	ErrorChan	 chan *Error
	client       *Client
}

// Sends request to unsubscribe
func (cs *CommandSubscription) Remove() *Error {
	commandSubsMutex.Lock()
	defer commandSubsMutex.Unlock()

	subsId := commandSubscriptions[cs]
	err := cs.client.unsubscribe("command/unsubscribe", subsId)

	if err != nil {
		return err
	}

	delete(commandSubscriptions, cs)

	return nil
}

func newCommandSubscription(subsId string, tspSubs *transport.Subscription, client *Client) *CommandSubscription {
	subs := &CommandSubscription{
		CommandsChan: make(chan *Command),
		ErrorChan:	  make(chan *Error),
		client:       client,
	}
	commandSubsMutex.Lock()
	commandSubscriptions[subs] = subsId
	commandSubsMutex.Unlock()

	go func() {
		loop: for {
			select {
			case rawComm, ok := <- tspSubs.DataChan:
				if !ok {
					break loop
				}

				comm := client.NewCommand()
				err := json.Unmarshal(rawComm, comm)
				if err != nil {
					subs.ErrorChan <- &Error{name: InvalidSubscriptionEventData, reason: err.Error()}
					continue
				}

				subs.CommandsChan <- comm
			case err, ok := <- tspSubs.ErrChan:
				if !ok {
					break loop
				}

				subs.ErrorChan <- newError(err)
			}
		}

		close(subs.CommandsChan)
		close(subs.ErrorChan)
	}()

	return subs
}
