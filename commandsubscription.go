package devicehive_go

import (
	"encoding/json"
	"log"
	"sync"
)

var commandSubsMutex = sync.Mutex{}
var commandSubscriptions = make(map[*CommandSubscription]string)

type CommandSubscription struct {
	CommandsChan chan *command
	client       *Client
}

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
		CommandsChan: make(chan *command),
		client:       client,
	}

	go func() {
		for rawComm := range tspChan {
			if client.transport.IsWS() {
				comm := &command{
					client: client,
				}
				err := json.Unmarshal(rawComm, &commandResponse{Command: comm})

				if err != nil {
					log.Println("couldn't unmarshal command data in subscription:", err)
					continue
				}

				subs.CommandsChan <- comm
			} else {
				var comms []*command
				err := json.Unmarshal(rawComm, &comms)
				for _, v := range comms {
					v.client = client
				}
				if err != nil {
					log.Println("couldn't unmarshal array of command data in subscription:", err)
					continue
				}

				for _, comm := range comms {
					subs.CommandsChan <- comm
				}
			}
		}

		close(subs.CommandsChan)
	}()

	commandSubsMutex.Lock()
	commandSubscriptions[subs] = subsId
	commandSubsMutex.Unlock()

	return subs
}
