package dh_test

import (
	"github.com/devicehive/devicehive-go/dh"
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestNotification(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is := is.New(t)

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	name := "test notif"
	ts := time.Now()
	params := map[string]interface{}{
		"testParam": 1,
	}
	id, dhErr := client.NotificationInsert(devId, name, ts, params)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(id != 0)

	notif, dhErr := client.NotificationGet(devId, id)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(notif != nil)
	is.Equal(int(notif.Parameters["testParam"].(float64)), 1)

	list, dhErr := client.NotificationList(devId, nil)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(len(list) > 0)
	is.Equal(int(list[0].Parameters["testParam"].(float64)), 1)

	listParams := &dh.ListParams{
		Start: time.Now().Add(-1 * time.Hour),
		End:   time.Now().Add(-1 * time.Minute),
	}
	list, dhErr = client.NotificationList(devId, listParams)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(len(list) == 0)
}

func TestNotificationSubscribe(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	is := is.New(t)

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	name := "test notif"
	ts := time.Now()

	notifChan, err := client.NotificationSubscribe(nil)

	go func() {
		select {
		case notif := <-notifChan:
			is.Equal(notif.Notification, name)
		case <-time.After(1 * time.Second):
			t.Error("notification insert event timeout")
		}
	}()

	_, err = client.NotificationInsert(devId, name, ts, nil)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	<-time.After(500 * time.Millisecond)
}

func TestNotificationUnsubscribe(t *testing.T) {
	err := auth()

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	devId := "4NemW3PE9BHRSqb0DVVgsphZh7SCZzgm3Lxg"
	name := "test notif"
	ts := time.Now()

	notifChan, err := client.NotificationSubscribe(nil)

	go func() {
		select {
		case notif, ok := <-notifChan:
			if notif != nil || ok {
				t.Error("client hasn't been unsubscribed")
			}
		case <-time.After(1 * time.Second):
			t.Error("timeout")
		}
	}()

	err = client.NotificationUnsubscribe(notifChan)

	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
		return
	}

	client.NotificationInsert(devId, name, ts, nil)
}
