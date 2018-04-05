package dh

import (
	"encoding/json"
)

type notification struct {
	Value *Notification `json:"notification"`
}

type Notification struct {
	Id int64 `json:"id"`
	Name string `json:"notification"`
	Timestamp dhTime `json:"timestamp"`
	DeviceId string `json:"deviceId"`
	NetworkId int64 `json:"networkId"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (c *Client) NotificationGet(deviceId string, notifId int64) (notif *Notification, err *Error) {
	_, rawRes, err := c.request(map[string]interface{}{
		"action": "notification/get",
		"deviceId": deviceId,
		"notificationId": notifId,
	})

	if err != nil {
		return nil, err
	}

	notif = &Notification{}
	parseErr := json.Unmarshal(rawRes, &notification{ Value: notif })

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return notif, nil
}
