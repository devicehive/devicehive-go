package devicehive_go

import (
	"encoding/json"
	"time"
)

type Device struct {
	Id           string                 `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`
	NetworkId    int                    `json:"networkId,omitempty"`
	DeviceTypeId int                    `json:"deviceTypeId,omitempty"`
	IsBlocked    bool                   `json:"isBlocked,omitempty"`
}

type device struct {
	Device
	client       *Client
}

func (d *device) Remove() *Error {
	_, err := d.client.request("deleteDevice", map[string]interface{}{
		"deviceId": d.Id,
	})

	return err
}

func (d *device) Save() *Error {
	_, err := d.client.request("putDevice", map[string]interface{}{
		"deviceId": d.Id,
		"device":   d,
	})

	return err
}

func (d *device) ListCommands(params *ListParams) (list []*command, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = d.Id

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := d.client.request("listCommands", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	for _, c := range list {
		c.client = d.client
	}

	return list, nil
}

func (d *device) SendCommand(name string, params map[string]interface{}, lifetime int, timestamp time.Time,
	status string, result map[string]interface{}) (comm *command, err *Error) {

	comm = d.client.NewCommand()

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

	rawRes, err := d.client.request("insertCommand", map[string]interface{}{
		"deviceId": d.Id,
		"command":  comm,
	})

	if err != nil {
		return nil, err
	}

	parseErr := json.Unmarshal(rawRes, comm)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	comm.DeviceId = d.Id

	return comm, nil
}

func (d *device) ListNotifications(params *ListParams) (list []*Notification, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = d.Id

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := d.client.request("listNotifications", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (d *device) SendNotification(name string, params map[string]interface{}, timestamp time.Time) (notif *Notification, err *Error) {
	notif = &Notification{
		Notification: name,
	}

	if params != nil {
		notif.Parameters = params
	}
	if timestamp.Unix() > 0 {
		notif.Timestamp = ISO8601Time{Time: timestamp}
	}

	rawRes, err := d.client.request("insertNotification", map[string]interface{}{
		"deviceId":     d.Id,
		"notification": notif,
	})

	if err != nil {
		return nil, err
	}

	pErr := json.Unmarshal(rawRes, notif)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return notif, nil
}

func (d *device) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return d.subscribeCommands(names, timestamp, false)
}

func (d *device) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return d.subscribeCommands(names, timestamp, true)
}

func (d *device) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		DeviceId:              d.Id,
	}

	return d.client.SubscribeCommands(params)

}

func (d *device) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	params := &SubscribeParams{
		Names:     names,
		Timestamp: timestamp,
		DeviceId:  d.Id,
	}

	return d.client.SubscribeNotifications(params)
}
