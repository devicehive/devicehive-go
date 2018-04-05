package dh

import (
	"encoding/json"
)

type notificationResponse struct {
	Value *Notification `json:"notification"`
	List *[]*Notification `json:"notifications"`
}

type Notification struct {
	Id int64 `json:"id"`
	Notification string `json:"notification"`
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
	parseErr := json.Unmarshal(rawRes, &notificationResponse{ Value: notif })

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return notif, nil
}

func (c *Client) NotificationList(deviceId string, params *ListParams) (list []*Notification, err *Error) {
	params.DeviceId = deviceId

	data := params.Map()
	data["action"] = "notification/list"

	_, rawRes, err := c.request(data)

	if err != nil {
		return nil, err
	}

	parseErr := json.Unmarshal(rawRes, &notificationResponse{ List: &list })

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}
