package devicehive_go

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/transportadapter"
	"time"
)

type WSClient struct {
	transportAdapter *transportadapter.WSAdapter
	DataChan         chan []byte
	ErrorChan        chan error
}

func (wsc *WSClient) unsubscribe(resourceName, subscriptionId string) {
	go func() {
		err := wsc.transportAdapter.Unsubscribe(resourceName, subscriptionId, Timeout)
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
		tspChan, subscriptionId, rawErr := wsc.transportAdapter.Subscribe(resourceName, 0, data)
		if rawErr != nil {
			wsc.ErrorChan <- newTransportErr(rawErr)
			return
		}

		res, _ := json.Marshal(map[string]string{
			"subscriptionId": subscriptionId,
		})

		wsc.DataChan <- res

		for b := range tspChan {
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
		resBytes, err := wsc.transportAdapter.Request(resourceName, data, Timeout)
		if err != nil {
			wsc.ErrorChan <- newError(err)
			return
		}

		wsc.DataChan <- resBytes
	}()

	return nil
}

func (wsc *WSClient) SubscribeNotifications(params *SubscribeParams) *Error {
	return wsc.subscribe("subscribeNotifications", params)
}

func (wsc *WSClient) UnsubscribeNotifications(subscriptionId string) {
	wsc.unsubscribe("notification/unsubscribe", subscriptionId)
}

func (wsc *WSClient) SubscribeCommands(params *SubscribeParams) *Error {
	return wsc.subscribe("subscribeCommands", params)
}

func (wsc *WSClient) UnsubscribeCommands(subscriptionId string) {
	wsc.unsubscribe("command/unsubscribe", subscriptionId)
}

func (wsc *WSClient) Authorize(accessToken string) *Error {
	return wsc.request("auth", map[string]interface{}{
		"token": accessToken,
	})
}

func (wsc *WSClient) PutDevice(device Device) *Error {
	if device.Name == "" {
		device.Name = device.Id
	}

	return wsc.request("putDevice", map[string]interface{}{
		"deviceId": device.Id,
		"device":   device,
	})
}

func (wsc *WSClient) GetDevice(deviceId string) *Error {
	return wsc.request("getDevice", map[string]interface{}{
		"deviceId": deviceId,
	})
}

func (wsc *WSClient) DeleteDevice(deviceId string) *Error {
	return wsc.request("deleteDevice", map[string]interface{}{
		"deviceId": deviceId,
	})
}

func (wsc *WSClient) UpdateDevice(deviceId string, device *Device) *Error {
	return wsc.request("putDevice", map[string]interface{}{
		"deviceId": deviceId,
		"device":   device,
	})
}

func (wsc *WSClient) ListDevices(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request("listDevices", data)
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

	return wsc.request("insertCommand", map[string]interface{}{
		"deviceId": deviceId,
		"command":  comm,
	})
}

func (wsc *WSClient) ListDeviceCommands(deviceId string, params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request("listCommands", data)
}

func (wsc *WSClient) SendDeviceNotification(deviceId, name string, params map[string]interface{}, timestamp time.Time) *Error {
	notif := &Notification{
		Notification: name,
	}

	notif.Parameters = params
	if timestamp.Unix() > 0 {
		notif.Timestamp = ISO8601Time{Time: timestamp}
	}

	return wsc.request("insertNotification", map[string]interface{}{
		"deviceId":     deviceId,
		"notification": notif,
	})
}

func (wsc *WSClient) ListDeviceNotifications(deviceId string, params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	params.DeviceId = deviceId

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request("listNotifications", data)
}

func (wsc *WSClient) GetInfo() *Error {
	return wsc.request("apiInfo", nil)
}

func (wsc *WSClient) GetClusterInfo() *Error {
	return wsc.request("apiInfoCluster", nil)
}

func (wsc *WSClient) SetProperty(name, value string) *Error {
	return wsc.request("putConfig", map[string]interface{}{
		"name":  name,
		"value": value,
	})
}

func (wsc *WSClient) GetProperty(name string) *Error {
	return wsc.request("getConfig", map[string]interface{}{
		"name": name,
	})
}

func (wsc *WSClient) DeleteProperty(name string) *Error {
	return wsc.request("deleteConfig", map[string]interface{}{
		"name": name,
	})
}

func (wsc *WSClient) CreateDeviceType(name, description string) *Error {

	devType := &DeviceType{
		Name:        name,
		Description: description,
	}

	return wsc.request("insertDeviceType", map[string]interface{}{
		"deviceType": devType,
	})

}

func (wsc *WSClient) GetDeviceType(deviceTypeId int) *Error {
	return wsc.request("getDeviceType", map[string]interface{}{
		"deviceTypeId": deviceTypeId,
	})

}

func (wsc *WSClient) ListDeviceTypes(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}
	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}
	return wsc.request("listDeviceTypes", data)

}

func (wsc *WSClient) CreateNetwork(name, description string) *Error {
	network := &Network{
		Name:        name,
		Description: description,
	}

	return wsc.request("insertNetwork", map[string]interface{}{
		"network": network,
	})
}

func (wsc *WSClient) GetNetwork(networkId int) *Error {

	return wsc.request("getNetwork", map[string]interface{}{
		"networkId": networkId,
	})

}

func (wsc *WSClient) ListNetworks(params *ListParams) *Error {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	return wsc.request("listNetworks", data)

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

	return wsc.request("tokenCreate", map[string]interface{}{
		"payload": data,
	})
}

func (wsc *WSClient) AccessTokenByRefresh(refreshToken string) *Error {
	return wsc.request("tokenRefresh", map[string]interface{}{
		"refreshToken": refreshToken,
	})
}

func (wsc *WSClient) AccessTokenByCreds(login, password string) *Error {
	return wsc.request("tokenByCreds", map[string]interface{}{
		"login":    login,
		"password": password,
	})
}
