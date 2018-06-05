package devicehive_go

type commandResponse struct {
	Command *Command `json:"command"`
}
type Command struct {
	Id          int         `json:"id,omitempty"`
	Command     string      `json:"command,omitempty"`
	Timestamp   ISO8601Time `json:"timestamp,omitempty"`
	LastUpdated ISO8601Time `json:"lastUpdated,omitempty"`
	UserId      int         `json:"userId,omitempty"`
	DeviceId    string      `json:"deviceId,omitempty"`
	NetworkId   int         `json:"networkId,omitempty"`
	Parameters  interface{} `json:"parameters,omitempty"`
	Lifetime    int         `json:"lifetime,omitempty"`
	Status      string      `json:"status,omitempty"`
	Result      interface{} `json:"result,omitempty"`
	client 		*Client
}

// Sends request to modify command at DeviceHive
func (comm *Command) Save() *Error {
	_, err := comm.client.request("updateCommand", map[string]interface{}{
		"deviceId":  comm.DeviceId,
		"commandId": comm.Id,
		"command":   comm,
	})

	return err
}
