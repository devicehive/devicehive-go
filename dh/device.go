package dh

import "encoding/json"

type deviceResponse struct {
	Device *Device `json:"device"`
}

type Device struct {
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
	NetworkId int64 `json:"networkId,omitempty"`
	DeviceTypeId int64 `json:"deviceTypeId,omitempty"`
	IsBlocked bool `json:"isBlocked,omitempty"`
	client *Client
}

func (c *Client) GetDevice(deviceId string) (device *Device, err *Error) {
	_, rawRes, err := c.request(map[string]interface{} {
		"action": "device/get",
		"deviceId": deviceId,
	})

	if err != nil {
		return nil, err
	}

	device = &Device{}
	parseErr := json.Unmarshal(rawRes, &deviceResponse{ Device: device })

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return device, nil
}

func (c *Client) PutDevice(deviceId string, device *Device) *Error {
	if device == nil {
		device = &Device{}
	}

	device.Id = deviceId

	if device.Name == "" {
		device.Name = deviceId
	}

	_, _, err := c.request(map[string]interface{} {
		"action": "device/save",
		"deviceId": deviceId,
		"device": device,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveDevice(deviceId string) *Error {
	_, _, err := c.request(map[string]interface{} {
		"action": "device/delete",
		"deviceId": deviceId,
	})

	return err
}
