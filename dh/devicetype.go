package dh

import (
	"encoding/json"
	"time"
	"strconv"
)

type DeviceType struct {
	client      *Client
	Id          int64  `json:"id,omitempty"`
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

func (c *Client) CreateDeviceType(name, description string) (devType *DeviceType, err *Error) {
	devType = &DeviceType{
		client:      c,
		Name:        name,
		Description: description,
	}

	res, err := c.request("insertDeviceType", map[string]interface{}{
		"deviceType": devType,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, devType)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return devType, nil
}

func (c *Client) GetDeviceType(deviceTypeId int64) (devType *DeviceType, err *Error) {
	devType = &DeviceType{
		client: c,
	}

	err = c.getModel("getDeviceType", devType, map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})
	if err != nil {
		return nil, err
	}

	return devType, nil
}

func (c *Client) ListDeviceTypes(params *ListParams) (list []*DeviceType, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listDeviceTypes", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (dt *DeviceType) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return dt.subscribeCommands(names, timestamp, false)
}

func (dt *DeviceType) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return dt.subscribeCommands(names, timestamp, true)
}

func (dt *DeviceType) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	id := []string{strconv.FormatInt(dt.Id, 10)}
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		DeviceTypeIds:         id,
	}

	return dt.client.SubscribeCommands(params)
}

func (dt *DeviceType) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	id := []string{strconv.FormatInt(dt.Id, 10)}
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		DeviceTypeIds:         id,
	}

	return dt.client.SubscribeNotifications(params)
}
