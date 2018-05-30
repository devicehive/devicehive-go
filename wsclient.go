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
