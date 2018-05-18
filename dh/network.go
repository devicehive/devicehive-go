package dh

import (
	"encoding/json"
)

type Network struct {
	client *Client
	Id int64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (n *Network) Save() *Error {
	_, err := n.client.request("updateNetwork", map[string]interface{}{
		"networkId": n.Id,
		"network": n,
	})

	return err
}

func (n *Network) Remove() *Error {
	_, err := n.client.request("deleteNetwork", map[string]interface{}{
		"networkId": n.Id,
	})

	return err
}

func (c *Client) CreateNetwork(name, description string) (network *Network, err *Error) {
	network = &Network{
		client: c,
		Name: name,
		Description: description,
	}

	res, err := c.request("insertNetwork", map[string]interface{}{
		"network": network,
	})
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(res, network)
	if jsonErr != nil {
		return nil, newJSONErr()
	}

	return network, nil
}

func (c *Client) GetNetwork(networkId int64) (network *Network, err *Error) {
	rawRes, err := c.request("getNetwork", map[string]interface{}{
		"networkId": networkId,
	})

	if err != nil {
		return nil, err
	}

	network = &Network{
		client: c,
	}
	parseErr := json.Unmarshal(rawRes, network)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	return network, nil
}
