package devicehive_go

import (
	"encoding/json"
	"strconv"
	"time"
)

type Network struct {
	client      *Client
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (n *Network) Save() *Error {
	_, err := n.client.request("updateNetwork", map[string]interface{}{
		"networkId": n.Id,
		"network":   n,
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
		client:      c,
		Name:        name,
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
	network = &Network{
		client: c,
	}

	err = c.getModel("getNetwork", network, map[string]interface{}{
		"networkId": networkId,
	})
	if err != nil {
		return nil, err
	}

	return network, nil
}

func (c *Client) ListNetworks(params *ListParams) (list []*Network, err *Error) {
	if params == nil {
		params = &ListParams{}
	}

	data, pErr := params.Map()
	if pErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: pErr.Error()}
	}

	rawRes, err := c.request("listNetworks", data)
	if err != nil {
		return nil, err
	}

	pErr = json.Unmarshal(rawRes, &list)
	if pErr != nil {
		return nil, newJSONErr()
	}

	return list, nil
}

func (n *Network) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return n.subscribeCommands(names, timestamp, false)
}

func (n *Network) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return n.subscribeCommands(names, timestamp, true)
}

func (n *Network) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	id := []string{strconv.FormatInt(n.Id, 10)}
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		NetworkIds:            id,
	}

	return n.client.SubscribeCommands(params)
}

func (n *Network) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	id := []string{strconv.FormatInt(n.Id, 10)}
	params := &SubscribeParams{
		Names:      names,
		Timestamp:  timestamp,
		NetworkIds: id,
	}

	return n.client.SubscribeNotifications(params)
}