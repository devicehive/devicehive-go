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

func (wsc *WSClient) SubscribeCommands(params *SubscribeParams) *Error {
	return wsc.subscribe("subscribeCommands", params)
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
		tspChan, _, rawErr := wsc.transportAdapter.Subscribe(resourceName, 0, data)
		if rawErr != nil {
			wsc.ErrorChan <- newTransportErr(rawErr)
			return
		}

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

func (wsc *WSClient) Authorize(accessToken string) *Error {
	return wsc.request("auth", map[string]interface{}{
		"token": accessToken,
	})
}

func (wsc *WSClient) PutDevice(deviceId, name string, data map[string]interface{}, networkId, deviceTypeId int, isBlocked bool) *Error {
	device := &Device{
		Id: deviceId,
	}

	if name == "" {
		device.Name = deviceId
	} else {
		device.Name = name
	}

	device.Data = data
	device.NetworkId = networkId
	device.DeviceTypeId = deviceTypeId
	device.IsBlocked = isBlocked

	return wsc.request("putDevice", map[string]interface{}{
		"deviceId": deviceId,
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
