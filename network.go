package devicehive_go

import (
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
