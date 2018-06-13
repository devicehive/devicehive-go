// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package stubs

var ResponseStub = &responseStub{
	actionRes: map[string]func(reqData map[string]interface{}) map[string]interface{}{
		"notification/subscribe": notificationSubscribe,
	},
}

type responseStub struct {
	actionRes map[string]func(reqData map[string]interface{}) map[string]interface{}
}

func (s *responseStub) Respond(reqData map[string]interface{}) map[string]interface{} {
	f, ok := s.actionRes[reqData["action"].(string)]

	if !ok {
		return notFound(reqData)
	}

	return f(reqData)
}

func notFound(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    reqData["action"],
		"requestId": reqData["requestId"],
		"status":    "error",
		"code":      404,
		"error":     "action not found",
	}
}

func (s *responseStub) NotificationInsertEvent(subscriptionId, deviceId interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":         "notification/insert",
		"subscriptionId": subscriptionId,
		"notification": map[string]interface{}{
			"id":           1,
			"notification": "notif test name",
			"timestamp":    "2018-04-03T05:57:59.379",
			"deviceId":     deviceId,
			"networkId":    1111,
			"parameters": map[string]interface{}{
				"testParam": 1,
			},
		},
	}
}

func notificationSubscribe(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":         "notification/subscribe",
		"status":         "success",
		"requestId":      reqData["requestId"],
		"subscriptionId": 1,
	}
}
