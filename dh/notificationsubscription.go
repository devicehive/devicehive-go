package dh

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
			if client.transport.IsWS() {
				notif := &Notification{}
				err := json.Unmarshal(rawNotif, &notificationResponse{Notification: notif})

				if err != nil {
					log.Println("couldn't unmarshal notification insert event data:", err)
					continue
				}

				subs.NotificationChan <- notif
			} else {
				var notifs []*Notification
				err := json.Unmarshal(rawNotif, &notifs)

				if err != nil {
					log.Println("couldn't unmarshal array of notification data in subscription:", err)
					continue
				}

				for _, notif := range notifs {
					subs.NotificationChan <- notif
				}
			}
		}

		close(subs.NotificationChan)
	}()

	notifSubsMutex.Lock()
	notificationSubscriptions[subs] = subsId
	notifSubsMutex.Unlock()

	return subs
}
