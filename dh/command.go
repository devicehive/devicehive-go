package dh

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
	client      *Client
}

func (comm *Command) Save() *Error {
	_, _, err := comm.client.request("command/update", map[string]interface{}{
		"deviceId":  comm.DeviceId,
		"commandId": comm.Id,
		"command":   comm,
	})

	return err
}
