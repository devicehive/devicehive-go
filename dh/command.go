package dh

import (
	"encoding/json"
)

type commandResponse struct {
	Command *Command `json:"command"`
	List *[]*Command `json:"commands"`
}

type Command struct {
	Id          int64                  `json:"id,omitempty"`
	Command     string                 `json:"command,omitempty"`
	Timestamp   ISO8601Time            `json:"timestamp,omitempty"`
	LastUpdated ISO8601Time            `json:"lastUpdated,omitempty"`
	UserId      int                    `json:"userId,omitempty"`
	DeviceId    string                 `json:"deviceId,omitempty"`
	NetworkId   int                    `json:"networkId,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Lifetime    int                    `json:"lifetime,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Result      map[string]interface{} `json:"result,omitempty"`
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
	comm.DeviceId = deviceId
	comm.Command = commandName

	_, rawRes, err := c.request(map[string]interface{} {
		"action": "command/insert",
		"deviceId": deviceId,
		"command": comm,
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
	_, _, err := c.request(map[string]interface{} {
		"action": "command/update",
		"deviceId": deviceId,
		"commandId": commandId,
		"command": comm,
	})

	return err
}