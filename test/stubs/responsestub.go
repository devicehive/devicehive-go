package stubs

var ResponseStub = &responseStub{
	actionRes: map[string]func(reqData map[string]interface{}) map[string]interface{}{
		"authenticate":             emptySuccessResponse,
		"token":                    token,
		"token/refresh":            tokenRefresh,
		"token/create":             token,
		"server/info":              serverInfo,
		"cluster/info":             clusterInfo,
		"subscription/list":        subscriptionList,
		"configuration/get":        configurationGet,
		"configuration/put":        configurationPut,
		"configuration/delete":     emptySuccessResponse,
		"notification/get":         notificationGet,
		"notification/list":        notificationList,
		"notification/insert":      notificationInsert,
		"notification/subscribe":   notificationSubscribe,
		"notification/unsubscribe": notificationUnsubscribe,
	},
}

type responseStub struct {
	actionRes map[string]func(reqData map[string]interface{}) map[string]interface{}
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

func emptySuccessResponse(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    reqData["action"],
		"requestId": reqData["requestId"],
		"status":    "success",
	}
}

func token(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":       "token",
		"requestId":    reqData["requestId"],
		"status":       "success",
		"accessToken":  "accTok",
		"refreshToken": "refTok",
	}
}

func tokenRefresh(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":      "token/refresh",
		"requestId":   reqData["requestId"],
		"status":      "success",
		"accessToken": "accTok",
	}
}

func unauthorized(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    reqData["action"],
		"requestId": reqData["requestId"],
		"status":    "error",
		"code":      401,
		"error":     "Unauthorized",
	}
}

func serverInfo(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "server/info",
		"requestId": reqData["requestId"],
		"status":    "success",
		"info": map[string]interface{}{
			"apiVersion":      "4.0.0",
			"serverTimestamp": "2018-04-03T05:57:59.379",
			"restServerUrl":   "https://dh.com/rest/api",
		},
	}
}

func clusterInfo(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "cluster/info",
		"requestId": reqData["requestId"],
		"status":    "success",
		"clusterInfo": map[string]interface{}{
			"bootstrap.servers": "localhost:1111",
			"zookeeper.connect": "localhost:2222",
		},
	}
}

func subscriptionList(reqData map[string]interface{}) map[string]interface{} {
	subscriptions := []map[string]interface{}{
		{
			"subscriptionId": 1,
			"type":           reqData["type"],
			"deviceId":       "d1",
			"networkIds":     []string{"n1", "n2"},
			"deviceTypeIds":  []string{"dt1", "dt2"},
			"names":          []string{"n1", "n2"},
			"timestamp":      "2018-04-03T05:57:59.379",
		},
		{
			"subscriptionId": 2,
			"type":           reqData["type"],
			"deviceId":       "d2",
			"networkIds":     nil,
			"deviceTypeIds":  nil,
			"names":          nil,
			"timestamp":      "2018-04-03T05:57:59.379",
		},
	}

	return map[string]interface{}{
		"action":        "subscription/list",
		"status":        "success",
		"requestId":     reqData["requestId"],
		"subscriptions": subscriptions,
	}
}

func configurationGet(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "configuration/get",
		"status":    "success",
		"requestId": reqData["requestId"],
		"configuration": map[string]interface{}{
			"name":          reqData["name"],
			"value":         "test value",
			"entityVersion": 2,
		},
	}
}

func configurationPut(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "configuration/put",
		"status":    "success",
		"requestId": reqData["requestId"],
		"configuration": map[string]interface{}{
			"name":          reqData["name"],
			"value":         reqData["value"],
			"entityVersion": 1,
		},
	}
}

func notificationGet(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "notification/get",
		"status":    "success",
		"requestId": reqData["requestId"],
		"notification": map[string]interface{}{
			"id":           reqData["notificationId"],
			"notification": "notif test name",
			"timestamp":    "2018-04-03T05:57:59.379",
			"deviceId":     reqData["deviceId"],
			"networkId":    1111,
			"parameters": map[string]interface{}{
				"testParam": 1,
			},
		},
	}
}

func notificationList(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "notification/list",
		"status":    "success",
		"requestId": reqData["requestId"],
		"notifications": []map[string]interface{}{
			{
				"id":           1111,
				"notification": "notif 1",
				"timestamp":    "2018-04-03T05:57:59.379",
				"deviceId":     reqData["deviceId"],
				"networkId":    1,
				"parameters": map[string]interface{}{
					"param1": 1,
				},
			},
			{
				"id":           2222,
				"notification": "notif 2",
				"timestamp":    "2018-04-03T06:57:59.379",
				"deviceId":     reqData["deviceId"],
				"networkId":    2,
				"parameters": map[string]interface{}{
					"param1": 2,
				},
			},
		},
	}
}

func notificationInsert(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "notification/list",
		"status":    "success",
		"requestId": reqData["requestId"],
		"notification": map[string]interface{}{
			"id":        1,
			"timestamp": "2018-04-03T05:57:59.379",
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

func notificationUnsubscribe(reqData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"action":    "notification/unsubscribe",
		"status":    "success",
		"requestId": reqData["requestId"],
	}
}
