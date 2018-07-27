// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"time"

	"github.com/devicehive/devicehive-go/internal/resourcenames"
	"github.com/devicehive/devicehive-go/internal/transportadapter"
)

type WSClient struct {
	transportAdapter *transportadapter.WSAdapter
	// Channel for receiving responses
	DataChan chan []byte
	// Channel for receiving errors
	ErrorChan             chan error
	defaultRequestTimeout time.Duration
}

func (wsc *WSClient) unsubscribe(resourceName, subscriptionId string) {
	go func() {
		err := wsc.transportAdapter.Unsubscribe(resourceName, subscriptionId, wsc.defaultRequestTimeout)
		if err != nil {
			wsc.ErrorChan <- newError(err)
		}
	}()
}

func (wsc *WSClient) subscribe(resourceName string, params *SubscribeParams) *Error {
	if params == nil {
		params = &SubscribeParams{}
	}

	data, jsonErr := params.Map()
	if jsonErr != nil {
		return &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	go func() {
		tspSubs, subscriptionId, rawErr := wsc.transportAdapter.Subscribe(resourceName, 0, data)
		if rawErr != nil {
			wsc.ErrorChan <- newTransportErr(rawErr)
			return
		}

		res, _ := json.Marshal(map[string]string{
			"subscriptionId": subscriptionId,
		})

		wsc.DataChan <- res

		for b := range tspSubs.DataChan {
			wsc.DataChan <- b
		}
	}()

	return nil
}

func (wsc *WSClient) request(resourceName string, data map[string]interface{}) *Error {
	_, err := json.Marshal(data)
	if err != nil {
		return &Error{name: InvalidRequestErr, reason: err.Error()}
	}

	go func() {
		resBytes, err := wsc.transportAdapter.Request(resourceName, data, wsc.defaultRequestTimeout)
		if err != nil {
			wsc.ErrorChan <- newError(err)
			return
		}

		wsc.DataChan <- resBytes
	}()

	return nil
}

// Subscribes for notifications with given params. If params is nil then default values take place.
// After successful subscription JSON object with only property "subscriptionId" is sent to the main data channel.
func (wsc *WSClient) SubscribeNotifications(params *SubscribeParams) *Error {
	return wsc.subscribe(resourcenames.SubscribeNotifications, params)
}

func (wsc *WSClient) UnsubscribeNotifications(subscriptionId string) {
	wsc.unsubscribe("notification/unsubscribe", subscriptionId)
}

// Subscribes for commands with given params. If params is nil then default values take place.
// After successful subscription JSON object with only property "subscriptionId" is sent to the main data channel.
func (wsc *WSClient) SubscribeCommands(params *SubscribeParams) *Error {
	return wsc.subscribe(resourcenames.SubscribeCommands, params)
}

func (wsc *WSClient) UnsubscribeCommands(subscriptionId string) {
	wsc.unsubscribe("command/unsubscribe", subscriptionId)
}

func (wsc *WSClient) Authenticate(accessToken string) *Error {
	return wsc.request(resourcenames.Auth, map[string]interface{}{
		"token": accessToken,
	})
}

func (wsc *WSClient) PutDevice(device Device) *Error {
	if device.Name == "" {
		device.Name = device.Id
	}

	return wsc.request(resourcenames.PutDevice, map[string]interface{}{
		"deviceId": device.Id,
		"device":   device,
	})
}

func (wsc *WSClient) GetDevice(deviceId string) *Error {
	return wsc.request(resourcenames.GetDevice, map[string]interface{}{
		"deviceId": deviceId,
	})
}

func (wsc *WSClient) DeleteDevice(deviceId string) *Error {
	return wsc.request(resourcenames.DeleteDevice, map[string]interface{}{
		"deviceId": deviceId,
	})
}

func (wsc *WSClient) UpdateDevice(deviceId string, device *Device) *Error {
	return wsc.request(resourcenames.PutDevice, map[string]interface{}{
		"deviceId": deviceId,
		"device":   device,
	})
}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListDevices(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request(resourcenames.ListDevices, data)
}

func (wsc *WSClient) SendDeviceCommand(deviceId, name string, params map[string]interface{}, lifetime int, timestamp time.Time,
	status string, result map[string]interface{}) *Error {

	comm := &Command{
		Command: name,
	}

	comm.Parameters = params
	comm.Lifetime = lifetime
	comm.Status = status
	comm.Result = result
	if timestamp.Unix() > 0 {
		comm.Timestamp = ISO8601Time{Time: timestamp}
	}

	return wsc.request(resourcenames.InsertCommand, map[string]interface{}{
		"deviceId": deviceId,
		"command":  comm,
	})
}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListDeviceCommands(deviceId string, params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request(resourcenames.ListCommands, data)
}

func (wsc *WSClient) SendDeviceNotification(deviceId, name string, params map[string]interface{}, timestamp time.Time) *Error {
	notif := &Notification{
		Notification: name,
	}

	notif.Parameters = params
	if timestamp.Unix() > 0 {
		notif.Timestamp = ISO8601Time{Time: timestamp}
	}

	return wsc.request(resourcenames.InsertNotification, map[string]interface{}{
		"deviceId":     deviceId,
		"notification": notif,
	})
}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListDeviceNotifications(deviceId string, params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request(resourcenames.ListNotifications, data)
}

func (wsc *WSClient) GetInfo() *Error {
	return wsc.request(resourcenames.ApiInfo, nil)
}

func (wsc *WSClient) GetClusterInfo() *Error {
	return wsc.request(resourcenames.ClusterInfo, nil)
}

func (wsc *WSClient) SetProperty(name, value string) *Error {
	return wsc.request(resourcenames.PutConfig, map[string]interface{}{
		"name":  name,
		"value": value,
	})
}

func (wsc *WSClient) GetProperty(name string) *Error {
	return wsc.request(resourcenames.GetConfig, map[string]interface{}{
		"name": name,
	})
}

func (wsc *WSClient) DeleteProperty(name string) *Error {
	return wsc.request(resourcenames.DeleteConfig, map[string]interface{}{
		"name": name,
	})
}

func (wsc *WSClient) CreateDeviceType(name, description string) *Error {

	devType := &DeviceType{
		Name:        name,
		Description: description,
	}

	return wsc.request(resourcenames.InsertDeviceType, map[string]interface{}{
		"deviceType": devType,
	})

}

func (wsc *WSClient) DeleteDeviceType(deviceTypeId int) *Error {
	return wsc.request(resourcenames.DeleteDeviceType, map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})
}

func (wsc *WSClient) GetDeviceType(deviceTypeId int) *Error {
	return wsc.request(resourcenames.GetDeviceType, map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})

}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListDeviceTypes(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}
	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}
	return wsc.request(resourcenames.ListDeviceTypes, data)

}

func (wsc *WSClient) CreateNetwork(name, description string) *Error {
	network := &Network{
		Name:        name,
		Description: description,
	}

	return wsc.request(resourcenames.InsertNetwork, map[string]interface{}{
		"network": network,
	})
}

func (wsc *WSClient) GetNetwork(networkId int) *Error {

	return wsc.request(resourcenames.GetNetwork, map[string]interface{}{
		"networkId": networkId,
	})

}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListNetworks(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request(resourcenames.ListNetworks, data)

}

func (wsc *WSClient) DeleteNetwork(networkId int) *Error {
	return wsc.request(resourcenames.DeleteNetwork, map[string]interface{}{
		"networkId": networkId,
	})
}

func (wsc *WSClient) CreateToken(userId int, expiration time.Time, actions, networkIds, deviceTypeIds []string) *Error {
	data := map[string]interface{}{
		"userId": userId,
	}

	if actions != nil {
		data["actions"] = actions
	}
	if networkIds != nil {
		data["networkIds"] = networkIds
	}
	if deviceTypeIds != nil {
		data["deviceTypeIds"] = deviceTypeIds
	}
	if expiration.Unix() > 0 {
		data["expiration"] = (&ISO8601Time{expiration}).String()
	}

	return wsc.request(resourcenames.TokenCreate, map[string]interface{}{
		"payload": data,
	})
}

func (wsc *WSClient) AccessTokenByRefresh(refreshToken string) *Error {
	return wsc.request(resourcenames.TokenRefresh, map[string]interface{}{
		"refreshToken": refreshToken,
	})
}

func (wsc *WSClient) AccessTokenByCreds(login, password string) *Error {
	return wsc.request(resourcenames.TokenByCreds, map[string]interface{}{
		"login":    login,
		"password": password,
	})
}

func (wsc *WSClient) CreateUser(login, password string, role int, data map[string]interface{}, allDevTypesAvail bool) *Error {
	return wsc.request(resourcenames.CreateUser, map[string]interface{}{
		"user": map[string]interface{}{
			"login":    login,
			"role":     role,
			"status":   UserStatusActive,
			"password": password,
			"data":     data,
			"allDeviceTypesAvailable": allDevTypesAvail,
		},
	})
}

func (wsc *WSClient) GetUser(userId int) *Error {
	return wsc.request(resourcenames.GetUser, map[string]interface{}{
		"userId": userId,
	})
}

func (wsc *WSClient) UpdateUser(userId int, user User) *Error {
	return wsc.request(resourcenames.UpdateUser, map[string]interface{}{
		"userId": userId,
		"user":   user,
	})
}

func (wsc *WSClient) DeleteUser(userId int) *Error {
	return wsc.request(resourcenames.DeleteUser, map[string]interface{}{
		"userId": userId,
	})
}

func (wsc *WSClient) GetCurrentUser() *Error {
	return wsc.request(resourcenames.GetCurrentUser, nil)
}

// In case params is nil default values defined at DeviceHive take place
func (wsc *WSClient) ListUsers(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request(resourcenames.ListUsers, data)
}

func (wsc *WSClient) UserAssignNetwork(userId, networkId int) *Error {
	return wsc.request(resourcenames.AssignNetwork, map[string]interface{}{
		"userId":    userId,
		"networkId": networkId,
	})
}

func (wsc *WSClient) UserUnassignNetwork(userId, networkId int) *Error {
	return wsc.request(resourcenames.UnassignNetwork, map[string]interface{}{
		"userId":    userId,
		"networkId": networkId,
	})
}

func (wsc *WSClient) UserAssignDeviceType(userId, deviceTypeId int) *Error {
	return wsc.request(resourcenames.AssignDeviceType, map[string]interface{}{
		"userId":       userId,
		"deviceTypeId": deviceTypeId,
	})
}

func (wsc *WSClient) UserUnassignDeviceType(userId, deviceTypeId int) *Error {
	return wsc.request(resourcenames.UnassignDeviceType, map[string]interface{}{
		"userId":       userId,
		"deviceTypeId": deviceTypeId,
	})
}

func (wsc *WSClient) AllowAllDeviceTypes(userId int) *Error {
	return wsc.request(resourcenames.AllowAllDeviceTypes, map[string]interface{}{
		"userId": userId,
	})
}

func (wsc *WSClient) DisallowAllDeviceTypes(userId int) *Error {
	return wsc.request(resourcenames.DisallowAllDeviceTypes, map[string]interface{}{
		"userId": userId,
	})
}

func (wsc *WSClient) ListUserDeviceTypes(userId int) *Error {
	return wsc.request(resourcenames.GetUserDeviceTypes, map[string]interface{}{
		"userId": userId,
	})
}
