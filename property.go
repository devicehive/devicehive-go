package devicehive_go

import (
	"encoding/json"
)

type Configuration struct {
	Name          string `json:"name"`
	Value         string `json:"value"`
	EntityVersion int    `json:"entityVersion"`
}

func (c *Client) GetProperty(name string) (conf *Configuration, err *Error) {
	conf = &Configuration{}

	err = c.getModel("getConfig", conf, map[string]interface{}{
		"name": name,
	})
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Client) SetProperty(name, value string) (entityVersion int, err *Error) {
	rawRes, err := c.request("putConfig", map[string]interface{}{
		"name":  name,
		"value": value,
	})

	if err != nil {
		return -1, err
	}

	conf := &Configuration{}
	parseErr := json.Unmarshal(rawRes, conf)
	if parseErr != nil {
		return -1, newJSONErr()
	}

	return conf.EntityVersion, nil
}

func (c *Client) DeleteProperty(name string) *Error {
	_, err := c.request("deleteConfig", map[string]interface{}{
		"name": name,
	})

	return err
}
