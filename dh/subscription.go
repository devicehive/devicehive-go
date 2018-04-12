package dh

import (
	"encoding/json"
)

type subscriptions struct {
	List []*Subscription `json:"subscriptions"`
}

type Subscription struct {
	Id            int64       `json:"subscriptionId"`
	Type          string      `json:"type"`
	DeviceId      string      `json:"deviceId"`
	NetworkIds    []string    `json:"networkIds"`
	DeviceTypeIds []string    `json:"deviceTypeIds"`
	Names         []string    `json:"names"`
	Timestamp     ISO8601Time `json:"timestamp"`
}

func (c *Client) SubscriptionList(subsType string) (list []*Subscription, err *Error) {
	_, resBytes, err := c.request(map[string]interface{}{
		"action": "subscription/list",
		"type":   subsType,
	})

	if err != nil {
		return nil, err
	}

	subs := &subscriptions{List: list}
	parseErr := json.Unmarshal(resBytes, subs)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return subs.List, nil
}
