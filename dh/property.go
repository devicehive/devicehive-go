package dh

import "encoding/json"

type configuration struct {
	Value *Configuration `json:"configuration"`
}

type Configuration struct {
	Name          string `json:"name"`
	Value         string `json:"value"`
	EntityVersion int    `json:"entityVersion"`
}

func (c *Client) GetProperty(name string) (conf *Configuration, err *Error) {
	rawRes, err := c.request("getConfig", map[string]interface{}{
		"name": name,
	})

	if err != nil {
		return nil, err
	}

	conf = &Configuration{}
	parseErr := json.Unmarshal(rawRes, &configuration{Value: conf})

	if parseErr != nil {
		return nil, newJSONErr()
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
	parseErr := json.Unmarshal(rawRes, &configuration{Value: conf})

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
