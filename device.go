// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

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
	client       *Client
}

func (d *Device) Remove() *Error {
	_, err := d.client.request("deleteDevice", map[string]interface{}{
		"deviceId": d.Id,
	})

	return err
}

func (d *Device) Save() *Error {
	_, err := d.client.request("putDevice", map[string]interface{}{
		"deviceId": d.Id,
		"device":   d,
	})

	return err
}

// In case params is nil default values defined at DeviceHive take place
func (d *Device) ListCommands(params *ListParams) ([]*Command, *Error) {
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

	var list []*Command
	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	for _, c := range list {
		c.client = d.client
	}

	return list, nil
}

func (d *Device) SendCommand(name string, params map[string]interface{}, lifetime int, timestamp time.Time,
	status string, result map[string]interface{}) (*Command, *Error) {

	comm := d.client.NewCommand()

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
	comm.Command = name

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

// In case params is nil default values defined at DeviceHive take place
func (d *Device) ListNotifications(params *ListParams) ([]*Notification, *Error) {
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

	var list []*Notification
	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (d *Device) SendNotification(name string, params map[string]interface{}, timestamp time.Time) (*Notification, *Error) {
	notif := &Notification{
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

func (d *Device) SubscribeInsertCommands(names []string, timestamp time.Time) (*CommandSubscription, *Error) {
	return d.subscribeCommands(names, timestamp, false)
}

func (d *Device) SubscribeUpdateCommands(names []string, timestamp time.Time) (*CommandSubscription, *Error) {
	return d.subscribeCommands(names, timestamp, true)
}

func (d *Device) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (*CommandSubscription, *Error) {
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		DeviceId:              d.Id,
	}

	return d.client.SubscribeCommands(params)

}

func (d *Device) SubscribeNotifications(names []string, timestamp time.Time) (*NotificationSubscription, *Error) {
	params := &SubscribeParams{
		Names:     names,
		Timestamp: timestamp,
		DeviceId:  d.Id,
	}

	return d.client.SubscribeNotifications(params)
}
