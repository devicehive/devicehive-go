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
	_, rawRes, err := c.request(map[string]interface{}{
		"action": "configuration/get",
		"name":   name,
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
	_, rawRes, err := c.request(map[string]interface{}{
		"action": "configuration/put",
		"name":   name,
		"value":  value,
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
	_, _, err := c.request(map[string]interface{}{
		"action": "configuration/delete",
		"name":   name,
	})

	return err
}