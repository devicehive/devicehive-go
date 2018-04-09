package dh_test

import (
	"testing"
	"time"
	"github.com/matryer/is"
	"github.com/devicehive/devicehive-go/dh"
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
	params := map[string]interface{} {
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
		End: time.Now().Add(-1 * time.Minute),
	}
	list, dhErr = client.NotificationList(devId, listParams)
	if dhErr != nil {
		t.Errorf("%s: %v", dhErr.Name(), dhErr)
		return
	}

	is.True(len(list) == 0)
}
