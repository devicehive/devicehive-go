// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"github.com/devicehive/devicehive-go/internal/resourcenames"
	"time"
)

type DeviceType struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	client      *Client
}

func (dt *DeviceType) Save() *Error {
	_, err := dt.client.request(resourcenames.UpdateDeviceType, map[string]interface{}{
		"deviceTypeId": dt.Id,
		"deviceType":   dt,
	})

	return err
}

func (dt *DeviceType) Remove() *Error {
	_, err := dt.client.request(resourcenames.DeleteDeviceType, map[string]interface{}{
		"deviceTypeId": dt.Id,
	})

	return err
}

func (dt *DeviceType) ForceRemove() *Error {
	_, err := dt.client.request(resourcenames.DeleteDeviceType, map[string]interface{}{
		"deviceTypeId": dt.Id,
		"force":        true,
	})

	return err
}

func (dt *DeviceType) SubscribeInsertCommands(names []string, timestamp time.Time) (*CommandSubscription, *Error) {
	return dt.subscribeCommands(names, timestamp, false)
}

func (dt *DeviceType) SubscribeUpdateCommands(names []string, timestamp time.Time) (*CommandSubscription, *Error) {
	return dt.subscribeCommands(names, timestamp, true)
}

func (dt *DeviceType) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (*CommandSubscription, *Error) {
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		DeviceTypeIds:         []int{dt.Id},
	}

	return dt.client.SubscribeCommands(params)
}

func (dt *DeviceType) SubscribeNotifications(names []string, timestamp time.Time) (*NotificationSubscription, *Error) {
	params := &SubscribeParams{
		Names:         names,
		Timestamp:     timestamp,
		DeviceTypeIds: []int{dt.Id},
	}

	return dt.client.SubscribeNotifications(params)
}
