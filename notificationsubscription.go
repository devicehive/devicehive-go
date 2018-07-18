// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/internal/transport"
	"sync"
)

var notifSubsMutex sync.RWMutex
var notificationSubscriptions = make(map[*NotificationSubscription]string)

type NotificationSubscription struct {
	NotificationChan chan *Notification
	ErrorChan        chan *Error
	done			 chan struct{}
	client           *Client
}

func (ns *NotificationSubscription) Remove() *Error {
	notifSubsMutex.RLock()
	subsId := notificationSubscriptions[ns]
	notifSubsMutex.RUnlock()

	err := ns.client.unsubscribe("notification/unsubscribe", subsId)
	if err != nil {
		return err
	}

	notifSubsMutex.Lock()
	close(ns.done)
	delete(notificationSubscriptions, ns)
	notifSubsMutex.Unlock()

	return nil
}

func (ns *NotificationSubscription) sendError(err *Error) {
	ns.ErrorChan <- err
}

func newNotificationSubscription(subsId string, tspSubs *transport.Subscription, client *Client) *NotificationSubscription {
	subs := &NotificationSubscription{
		NotificationChan: make(chan *Notification),
		ErrorChan:        make(chan *Error),
		done:		  	  make(chan struct{}),
		client:           client,
	}

	notifSubsMutex.Lock()
	notificationSubscriptions[subs] = subsId
	notifSubsMutex.Unlock()

	go func() {
	loop:
		for {
			select {
			case rawNotif, ok := <-tspSubs.DataChan:
				if !ok {
					break loop
				}

				notif := client.NewNotification()
				err := json.Unmarshal(rawNotif, notif)
				if err != nil {
					subs.ErrorChan <- &Error{name: InvalidSubscriptionEventData, reason: err.Error()}
					continue
				}

				subs.NotificationChan <- notif
			case err, ok := <-tspSubs.ErrChan:
				if !ok {
					break loop
				}

				client.handleSubscriptionError(subs, err)
			case <- subs.done:
				break loop
			}
		}

		close(subs.NotificationChan)
		close(subs.ErrorChan)
	}()

	return subs
}
