package stubs

var ResponseStub = &responseStub{}

type responseStub struct{}

func (s *responseStub) EmptySuccessResponse(action, reqId string) map[string]interface{} {
	return map[string]interface{}{
		"action":    action,
		"requestId": reqId,
		"status":    "success",
	}
}

func (s *responseStub) Token(reqId, accessToken, refreshToken string) map[string]interface{} {
	return map[string]interface{}{
		"action":       "token",
		"requestId":    reqId,
		"status":       "success",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
}

func (s *responseStub) TokenRefresh(reqId, accessToken string) map[string]interface{} {
	return map[string]interface{}{
		"action":      "token/refresh",
		"requestId":   reqId,
		"status":      "success",
		"accessToken": accessToken,
	}
}

func (s *responseStub) Unauthorized(action, reqId string) map[string]interface{} {
	return map[string]interface{}{
		"action":    action,
		"requestId": reqId,
		"status":    "error",
		"code":      401,
		"error":     "Unauthorized",
	}
}

func (s *responseStub) ServerInfo(reqId string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "server/info",
		"requestId": reqId,
		"status":    "success",
		"info": map[string]interface{}{
			"apiVersion":      "4.0.0",
			"serverTimestamp": "2018-04-03T05:57:59.379",
			"restServerUrl":   "https://dh.com/rest/api",
		},
	}
}

func (s *responseStub) ClusterInfo(reqId string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "cluster/info",
		"requestId": reqId,
		"status":    "success",
		"clusterInfo": map[string]interface{}{
			"bootstrap.servers": "localhost:1111",
			"zookeeper.connect": "localhost:2222",
		},
	}
}

func (s *responseStub) SubscriptionList(reqId, subsType string) map[string]interface{} {
	subscriptions := []map[string]interface{}{
		{
			"subscriptionId": 1,
			"type":           subsType,
			"deviceId":       "d1",
			"networkIds":     []string{"n1", "n2"},
			"deviceTypeIds":  []string{"dt1", "dt2"},
			"names":          []string{"n1", "n2"},
			"timestamp":      "2018-04-03T05:57:59.379",
		},
		{
			"subscriptionId": 2,
			"type":           subsType,
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
		"requestId":     reqId,
		"subscriptions": subscriptions,
	}
}

func (s *responseStub) ConfigurationGet(reqId, name string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "configuration/get",
		"status":    "success",
		"requestId": reqId,
		"configuration": map[string]interface{}{
			"name":          name,
			"value":         "test value",
			"entityVersion": 2,
		},
	}
}

func (s *responseStub) ConfigurationPut(reqId, name, val string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "configuration/put",
		"status":    "success",
		"requestId": reqId,
		"configuration": map[string]interface{}{
			"name":          name,
			"value":         val,
			"entityVersion": 1,
		},
	}
}
