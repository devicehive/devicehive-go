package dh

import (
	"fmt"
	"github.com/devicehive/devicehive-go/internal/utils"
	"time"
	"encoding/json"
)

type subscriptions struct {
	List []*Subscription `json:"subscription"`
}

type Subscription struct {
	Id            int64 `json:"subscriptionId"`
	Type          string `json:"type"`
	DeviceId      string `json:"deviceId"`
	NetworkIds    []string `json:"networkIds"`
	DeviceTypeIds []string `json:"deviceTypeIds"`
	Names         []string `json:"names"`
	Timestamp     time.Time `json:"timestamp"`
}

func (c *Client) SubscriptionList(subsType string) (list []*Subscription, err *Error) {
	resBytes, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "subscription/list",
		"type":   subsType,
	})

	if _, err = c.handleResponse(resBytes, tspErr); err != nil {
		return nil, err
	}

	subs := &subscriptions{ List: list }
	parseErr := json.Unmarshal(resBytes, subs)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	//for _, sub := range subs {
	//	rawSub := sub.(map[string]interface{})
	//	ts, tserr := time.Parse(timestampLayout, rawSub["timestamp"].(string))
	//
	//	if tserr != nil {
	//		return nil, &Error{name: InvalidResponseErr, reason: tserr.Error()}
	//	}
	//
	//	lists, err := normalizeLists(rawSub)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	list = append(list, &Subscription{
	//		Id:            int64(rawSub["subscriptionId"].(float64)),
	//		Type:          rawSub["type"].(string),
	//		DeviceId:      rawSub["deviceId"].(string),
	//		NetworkIds:    lists["networkIds"],
	//		DeviceTypeIds: lists["deviceTypeIds"],
	//		Names:         lists["names"],
	//		Timestamp:     ts,
	//	})
	//}

	return list, nil
}

func normalizeLists(rawSub map[string]interface{}) (map[string][]string, *Error) {
	listKeys := []string{"networkIds", "deviceTypeIds", "names"}
	strSlices := make(map[string][]string)

	for _, k := range listKeys {
		if rawSub[k] == nil {
			strSlices[k] = nil
			continue
		}

		res, err := utils.ISliceToStrSlice(rawSub[k].([]interface{}))

		if err != nil {
			r := fmt.Sprintf("%s is not array of strings", k)
			return nil, &Error{name: InvalidResponseErr, reason: r}
		}

		strSlices[k] = res
	}

	return strSlices, nil
}
