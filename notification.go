package devicehive_go

type notificationResponse struct {
	Notification *Notification    `json:"notification"`
	List         *[]*Notification `json:"notifications"`
}

type Notification struct {
	Id           int64                  `json:"id"`
	Notification string                 `json:"notification"`
	Timestamp    ISO8601Time            `json:"timestamp"`
	DeviceId     string                 `json:"deviceId"`
	NetworkId    int64                  `json:"networkId"`
	Parameters   map[string]interface{} `json:"parameters"`
}
