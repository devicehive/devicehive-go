package dh

import (
	"encoding/json"
	"time"
)

type deviceResponse struct {
	Device *Device `json:"device"`
}

type Device struct {
	Id           string                 `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	NetworkId    int64                  `json:"networkId,omitempty"`
	DeviceTypeId int64                  `json:"deviceTypeId,omitempty"`
	IsBlocked    bool                   `json:"isBlocked,omitempty"`
	client       *Client
}

func (d *Device) Remove() *Error {
	_, _, err := d.client.request("device/delete", map[string]interface{}{
		"deviceId": d.Id,
	})

	return err
}

func (d *Device) Save() *Error {
	_, _, err := d.client.request("device/save", map[string]interface{}{
		"deviceId": d.Id,
		"device":   d,
	})

	return err
}

func (d *Device) ListCommands(params *ListParams) (list []*Command, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = d.Id

	data, pErr := params.Map()

	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	_, rawRes, err := d.client.request("command/list", data)

	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &commandResponse{List: &list})

	if pErr != nil {
		return nil, newJSONErr()
	}

	for _, c := range list {
		c.client = d.client
	}

	return list, nil
}

func (d *Device) SendCommand(name string, params map[string]interface{}, lifetime int, timestamp time.Time,
	status string, result map[string]interface{}) (comm *Command, err *Error) {

	comm = &Command{
		Command: name,
		client:  d.client,
	}

	if params != nil {
		comm.Parameters = params
	}
	if lifetime != 0 {
		comm.Lifetime = lifetime
	}
	if timestamp.Unix() > 0 {
		comm.Timestamp = ISO8601Time{Time: timestamp}
	}
	if status != "" {
		comm.Status = status
	}
	if result != nil {
		comm.Result = result
	}

	_, rawRes, err := d.client.request("command/insert", map[string]interface{}{
		"deviceId": d.Id,
		"command":  comm,
	})

	if err != nil {
		return nil, err
	}

	parseErr := json.Unmarshal(rawRes, &commandResponse{Command: comm})

	if parseErr != nil {
		return nil, newJSONErr()
	}

	comm.DeviceId = d.Id

	return comm, nil
}

func (d *Device) ListNotifications(params *ListParams) (list []*Notification, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = d.Id

	data, pErr := params.Map()

	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	_, rawRes, err := d.client.request("notification/list", data)

	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &notificationResponse{List: &list})

	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (d *Device) SendNotification(name string, params map[string]interface{}, timestamp time.Time) (notif *Notification, err *Error) {
	notif = &Notification{
		Notification: name,
	}

	if params != nil {
		notif.Parameters = params
	}
	if timestamp.Unix() > 0 {
		notif.Timestamp = ISO8601Time{Time: timestamp}
	}

	_, rawRes, err := d.client.request("notification/insert", map[string]interface{}{
		"deviceId":     d.Id,
		"notification": notif,
	})

	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, &notificationResponse{Notification: notif})

	if pErr != nil {
		return nil, newJSONErr()
	}

	return notif, nil
}

func (d *Device) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return d.subscribeCommands(names, timestamp, false)
}

func (d *Device) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return d.subscribeCommands(names, timestamp, true)
}

func (d *Device) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	s, err := d.subscribe(&SubscribeParams{
		Names: names,
		Timestamp: timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
	}, "command/subscribe")

	if err != nil || s == nil {
		return nil, err
	}

	return s.(*CommandSubscription), nil

}

func (d *Device) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	s, err := d.subscribe(&SubscribeParams{
		Names: names,
		Timestamp: timestamp,
	}, "notification/subscribe")

	if err != nil || s == nil {
		return nil, err
	}

	return s.(*NotificationSubscription), nil
}

func (d *Device) subscribe(params *SubscribeParams, action string) (subs interface{}, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	params.DeviceId = d.Id

	tspChan, subsId, err := d.client.subscribe(action, params)

	if err != nil || tspChan == nil {
		return nil, err
	}

	if action == "notification/subscribe" {
		subs = newNotificationSubscription(subsId, tspChan, d.client)
	} else if action == "command/subscribe" {
		subs = newCommandSubscription(subsId, tspChan, d.client)
	}

	return subs, nil
}

func (c *Client) GetDevice(deviceId string) (device *Device, err *Error) {
	_, rawRes, err := c.request("device/get", map[string]interface{}{
		"deviceId": deviceId,
	})

	if err != nil {
		return nil, err
	}

	device = &Device{
		client: c,
	}
	parseErr := json.Unmarshal(rawRes, &deviceResponse{Device: device})

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return device, nil
}

func (c *Client) PutDevice(deviceId, name string, data map[string]interface{}, networkId, deviceTypeId int64, isBlocked bool) (device *Device, err *Error) {
	device = &Device{
		client: c,
	}

	device.Id = deviceId

	if name == "" {
		device.Name = deviceId
	} else {
		device.Name = name
	}

	if data != nil {
		device.Data = data
	}

	if networkId != 0 {
		device.NetworkId = networkId
	}

	if deviceTypeId != 0 {
		device.DeviceTypeId = deviceTypeId
	}

	if isBlocked {
		device.IsBlocked = isBlocked
	}

	_, _, err = c.request("device/save", map[string]interface{}{
		"deviceId": deviceId,
		"device":   device,
	})

	if err != nil {
		return nil, err
	}

	return device, nil
}
