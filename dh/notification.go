package dh

import (
	"encoding/json"
	"time"
)

var notificationSubscriptions = make(map[chan *Notification]string)

type notificationResponse struct {
	Notification *Notification    `json:"notification"`
	List         *[]*Notification `json:"notifications"`
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
	pErr := json.Unmarshal(rawRes, &notificationResponse{ Notification: notif })

	if pErr != nil {
		return nil, newJSONErr()
	}

	return notif, nil
}

func (c *Client) NotificationList(deviceId string, params *ListParams) (list []*Notification, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId
	params.Action = "notification/list"

	data, pErr := params.Map()

	if pErr != nil {
		return nil, &Error{ name: InvalidRequestErr, reason: pErr.Error() }
	}

	_, rawRes, err := c.request(data)

	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &notificationResponse{ List: &list })

	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (c *Client) NotificationInsert(deviceId, notifName string, timestamp time.Time, params map[string]interface{}) (notifId int64, err *Error) {
	_, rawRes, err := c.request(map[string]interface{} {
		"action": "notification/insert",
		"deviceId": deviceId,
		"notification": map[string]interface{} {
			"notification": notifName,
			"timestamp": timestamp.UTC().Format(timestampLayout),
			"parameters": params,
		},
	})

	if err != nil {
		return 0, err
	}

	notif := &Notification{}
	pErr := json.Unmarshal(rawRes, &notificationResponse{ Notification: notif })

	if pErr != nil {
		return 0, newJSONErr()
	}

	return notif.Id, nil
}

func (c *Client) NotificationSubscribe(params *SubscribeParams) (notifChan chan *Notification, err *Error) {
	tspChan, subsId,  err := c.subscribe("notification/subscribe", params)

	if err != nil {
		return nil, err
	}

	notifChan = c.notificationsTransform(tspChan)

	notificationSubscriptions[notifChan] = subsId

	return notifChan, nil
}

func (c *Client) NotificationUnsubscribe(notifChan chan *Notification) *Error {
	subsId := notificationSubscriptions[notifChan]
	return c.unsubscribe("notification/unsubscribe", subsId)
}

func (c *Client) notificationsTransform(tspChan chan []byte) (notifChan chan *Notification) {
	notifChan = make(chan *Notification)

	go func() {
		for rawNotif := range tspChan {
			notif := &Notification{}
			err := json.Unmarshal(rawNotif, &notificationResponse{ Notification: notif })

			if err != nil {
				close(notifChan)
				return
			}

			notifChan <- notif
		}

		close(notifChan)
	}()

	return notifChan
}
