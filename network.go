package devicehive_go

import (
	"time"
)

type network struct {
	client      *Client
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (n *network) Save() *Error {
	_, err := n.client.request("updateNetwork", map[string]interface{}{
		"networkId": n.Id,
		"network":   n,
	})

	return err
}

func (n *network) Remove() *Error {
	_, err := n.client.request("deleteNetwork", map[string]interface{}{
		"networkId": n.Id,
	})

	return err
}

func (n *network) ForceRemove() *Error {
	_, err := n.client.request("deleteNetwork", map[string]interface{}{
		"networkId": n.Id,
		"force":     true,
	})

	return err
}

func (n *network) SubscribeInsertCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return n.subscribeCommands(names, timestamp, false)
}

func (n *network) SubscribeUpdateCommands(names []string, timestamp time.Time) (subs *CommandSubscription, err *Error) {
	return n.subscribeCommands(names, timestamp, true)
}

func (n *network) subscribeCommands(names []string, timestamp time.Time, isCommUpdatesSubscription bool) (subs *CommandSubscription, err *Error) {
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		NetworkIds:            []int{n.Id},
	}

	return n.client.SubscribeCommands(params)
}

func (n *network) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	params := &SubscribeParams{
		Names:      names,
		Timestamp:  timestamp,
		NetworkIds: []int{n.Id},
	}

	return n.client.SubscribeNotifications(params)
}
