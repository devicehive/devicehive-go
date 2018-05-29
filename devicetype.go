package devicehive_go

import (
	"time"
)

type DeviceType struct {
	client      *Client
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (dt *DeviceType) Save() *Error {
	_, err := dt.client.request("updateDeviceType", map[string]interface{}{
		"deviceTypeId": dt.Id,
		"deviceType":   dt,
	})

	return err
}

func (dt *DeviceType) Remove() *Error {
	_, err := dt.client.request("deleteDeviceType", map[string]interface{}{
		"deviceTypeId": dt.Id,
	})

	return err
}

func (dt *DeviceType) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return dt.subscribeCommands(names, timestamp, false)
}

func (dt *DeviceType) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return dt.subscribeCommands(names, timestamp, true)
}

func (dt *DeviceType) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		DeviceTypeIds:         []int{dt.Id},
	}

	return dt.client.SubscribeCommands(params)
}

func (dt *DeviceType) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	params := &SubscribeParams{
		Names:         names,
		Timestamp:     timestamp,
		DeviceTypeIds: []int{dt.Id},
	}

	return dt.client.SubscribeNotifications(params)
}
