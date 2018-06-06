// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"time"
)

type Network struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	client 		*Client
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

func (n *Network) ForceRemove() *Error {
	_, err := n.client.request("deleteNetwork", map[string]interface{}{
		"networkId": n.Id,
		"force":     true,
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
	params := &SubscribeParams{
		Names:                 names,
		Timestamp:             timestamp,
		ReturnUpdatedCommands: isCommUpdatesSubscription,
		NetworkIds:            []int{n.Id},
	}

	return n.client.SubscribeCommands(params)
}

func (n *Network) SubscribeNotifications(names []string, timestamp time.Time) (subs *NotificationSubscription, err *Error) {
	params := &SubscribeParams{
		Names:      names,
		Timestamp:  timestamp,
		NetworkIds: []int{n.Id},
	}

	return n.client.SubscribeNotifications(params)
}
