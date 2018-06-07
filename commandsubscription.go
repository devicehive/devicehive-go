// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"log"
	"sync"
)

var commandSubsMutex = sync.Mutex{}
var commandSubscriptions = make(map[*CommandSubscription]string)

type CommandSubscription struct {
	CommandsChan chan *Command
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

func newCommandSubscription(subsId string, tspChan chan []byte, client *Client) *CommandSubscription {
	subs := &CommandSubscription{
		CommandsChan: make(chan *Command),
		client:       client,
	}
	commandSubsMutex.Lock()
	commandSubscriptions[subs] = subsId
	commandSubsMutex.Unlock()

	go func() {
		for rawComm := range tspChan {
			comm := client.NewCommand()
			err := json.Unmarshal(rawComm, comm)
			if err != nil {
				log.Printf("Error while parsing command subscription event: %s %s\n", err, string(rawComm))
				continue
			}

			subs.CommandsChan <- comm
		}

		close(subs.CommandsChan)
	}()

	return subs
}
