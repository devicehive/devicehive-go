package dh

import (
	"encoding/json"
	"log"
	"sync"
)

var commandSubsMutex = sync.Mutex{}
var commandSubscriptions = make(map[chan *Command]string)

type commandResponse struct {
	Command *Command    `json:"command"`
	List    *[]*Command `json:"commands"`
}

type Command struct {
	Id          int64                  `json:"id,omitempty"`
	Command     string                 `json:"command,omitempty"`
	Timestamp   ISO8601Time            `json:"timestamp,omitempty"`
	LastUpdated ISO8601Time            `json:"lastUpdated,omitempty"`
	UserId      int64                  `json:"userId,omitempty"`
	DeviceId    string                 `json:"deviceId,omitempty"`
	NetworkId   int64                  `json:"networkId,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Lifetime    int                    `json:"lifetime,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Result      map[string]interface{} `json:"result,omitempty"`
	client *Client
}

func (comm *Command) Save() *Error {
	_, _, err := comm.client.request(map[string]interface{}{
		"action":    "command/update",
		"deviceId":  comm.DeviceId,
		"commandId": comm.Id,
		"command":   comm,
	})

	return err
}

func (c *Client) CommandGet(deviceId string, commandId int64) (comm *Command, err *Error) {
	_, rawRes, err := c.request(map[string]interface{}{
		"action":    "command/get",
		"deviceId":  deviceId,
		"commandId": commandId,
	})

	if err != nil {
		return nil, err
	}

	comm = &Command{}
	pErr := json.Unmarshal(rawRes, &commandResponse{Command: comm})

	if pErr != nil {
		return nil, newJSONErr()
	}

	return comm, nil
}

func (c *Client) CommandList(deviceId string, params *ListParams) (list []*Command, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId
	params.Action = "command/list"

	data, pErr := params.Map()

	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	_, rawRes, err := c.request(data)

	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &commandResponse{List: &list})

	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (c *Client) CommandInsert(deviceId, commandName string, comm *Command) *Error {
	if comm == nil {
		comm = &Command{}
	}

	comm.DeviceId = deviceId
	comm.Command = commandName

	_, rawRes, err := c.request(map[string]interface{}{
		"action":   "command/insert",
		"deviceId": deviceId,
		"command":  comm,
	})

	if err != nil {
		return err
	}

	parseErr := json.Unmarshal(rawRes, &commandResponse{Command: comm})

	if parseErr != nil {
		return newJSONErr()
	}

	return nil
}

func (c *Client) CommandUpdate(deviceId string, commandId int64, comm *Command) *Error {
	_, _, err := c.request(map[string]interface{}{
		"action":    "command/update",
		"deviceId":  deviceId,
		"commandId": commandId,
		"command":   comm,
	})

	return err
}

func (c *Client) CommandSubscribe(params *SubscribeParams) (commChan chan *Command, err *Error) {
	tspChan, subsId, err := c.subscribe("command/subscribe", params)

	if err != nil {
		return nil, err
	}

	if tspChan == nil {
		return nil, nil
	}

	commChan = c.commandsTransform(tspChan)

	commandSubsMutex.Lock()
	commandSubscriptions[commChan] = subsId
	commandSubsMutex.Unlock()

	return commChan, nil
}

func (c *Client) CommandUnsubscribe(commandChan chan *Command) *Error {
	commandSubsMutex.Lock()
	defer commandSubsMutex.Unlock()

	subsId := commandSubscriptions[commandChan]
	err := c.unsubscribe("command/unsubscribe", subsId)

	if err != nil {
		return err
	}

	delete(commandSubscriptions, commandChan)

	return nil
}

func (c *Client) commandsTransform(tspChan chan []byte) (commChan chan *Command) {
	commChan = make(chan *Command)

	go func() {
		for rawComm := range tspChan {
			comm := &Command{}
			err := json.Unmarshal(rawComm, &commandResponse{Command: comm})

			if err != nil {
				log.Println("couldn't unmarshal command insert event data:", err)
				continue
			}

			commChan <- comm
		}

		close(commChan)
	}()

	return commChan
}
