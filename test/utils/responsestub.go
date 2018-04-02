package utils

var ResponseStub = &responseStub{}

type responseStub struct{}

func (s *responseStub) Authenticate(reqId string) map[string]interface{} {
	return map[string]interface{}{
		"action":    "authenticate",
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
			"serverTimestamp": "2006-01-02T15:04:05.000",
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
			"timestamp":      "2006-01-02T15:04:05.000",
		},
		{
			"subscriptionId": 2,
			"type":           subsType,
			"deviceId":       "d2",
			"networkIds":     nil,
			"deviceTypeIds":  nil,
			"names":          nil,
			"timestamp":      "2006-01-02T15:04:05.000",
		},
	}

	return map[string]interface{}{
		"action":        "subscription/list",
		"status":        "success",
		"requestId":     reqId,
		"subscriptions": subscriptions,
	}
}
