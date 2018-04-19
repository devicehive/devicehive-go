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

func (d *Device) Remove() *Error {
	_, _, err := d.client.request(map[string]interface{} {
		"action": "device/delete",
		"deviceId": d.Id,
	})

	return err
}

func (c *Client) GetDevice(deviceId string) (device *Device, err *Error) {
	_, rawRes, err := c.request(map[string]interface{} {
		"action": "device/get",
		"deviceId": deviceId,
	})

	if err != nil {
		return nil, err
	}

	device = &Device{
		client: c,
	}
	parseErr := json.Unmarshal(rawRes, &deviceResponse{ Device: device })

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return device, nil
}

func (c *Client) PutDevice(deviceId, name string, data map[string]interface{}, networkId, deviceTypeId int64, isBlocked bool) (device *Device, err *Error) {
	device = &Device{
		client: c,
	}

	device.Id = deviceId

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

	_, _, err = c.request(map[string]interface{} {
		"action": "device/save",
		"deviceId": deviceId,
		"device": device,
	})

	if err != nil {
		return nil, err
	}

	return device, nil
}
