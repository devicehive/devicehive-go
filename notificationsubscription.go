// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"encoding/json"
	"log"
	"sync"
)

var notifSubsMutex = sync.Mutex{}
var notificationSubscriptions = make(map[*NotificationSubscription]string)

type NotificationSubscription struct {
	NotificationChan chan *Notification
	client           *Client
}

func (ns *NotificationSubscription) Remove() *Error {
	notifSubsMutex.Lock()
	defer notifSubsMutex.Unlock()

	subsId := notificationSubscriptions[ns]
	err := ns.client.unsubscribe("notification/unsubscribe", subsId)

	if err != nil {
		return err
	}

	delete(notificationSubscriptions, ns)

	return nil
}

func newNotificationSubscription(subsId string, tspChan chan []byte, client *Client) *NotificationSubscription {
	subs := &NotificationSubscription{
		NotificationChan: make(chan *Notification),
		client:           client,
	}

	go func() {
		for rawNotif := range tspChan {
			notif := client.NewNotification()
			err := json.Unmarshal(rawNotif, notif)
			if err != nil {
				log.Printf("Error while parsing notification subscription event: %s %s\n", err, string(rawNotif))
				continue
			}

			subs.NotificationChan <- notif
		}

		close(subs.NotificationChan)
	}()

	notifSubsMutex.Lock()
	notificationSubscriptions[subs] = subsId
	notifSubsMutex.Unlock()

	return subs
}
