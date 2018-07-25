package responsehandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type wsResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

func WSHandleResponseError(rawRes []byte) error {
	res := &wsResponse{}
	parseErr := json.Unmarshal(rawRes, res)
	if parseErr != nil {
		return parseErr
	}

	if res.Status == "error" {
		errMsg := strings.ToLower(res.Error)
		errCode := res.Code
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return errors.New(r)
	}

	return nil
}

func WSExtractResponsePayload(resourceName string, rawRes []byte) []byte {
	payloadKey := wsResponsePayloads[resourceName]
	if payloadKey == "" {
		return rawRes
	}

	res := make(map[string]json.RawMessage)
	json.Unmarshal(rawRes, &res)

	return res[payloadKey]
}

var wsResponsePayloads = map[string]string{
	"getConfig":                   "configuration",
	"putConfig":                   "configuration",
	"deleteConfig":                "configuration",
	"apiInfo":                     "info",
	"apiInfoCluster":              "clusterInfo",
	"listCommands":                "commands",
	"insertCommand":               "command",
	"listNotifications":           "notifications",
	"insertNotification":          "notification",
	"subscribeNotificationsEvent": "notification",
	"subscribeCommandsEvent":      "command",
	"getDevice":                   "device",
	"commandEvent":                "command",
	"notificationEvent":           "notification",
	"listDevices":                 "devices",
	"insertNetwork":               "network",
	"getNetwork":                  "network",
	"listNetworks":                "networks",
	"insertDeviceType":            "deviceType",
	"getDeviceType":               "deviceType",
	"listDeviceTypes":             "deviceTypes",
	"createUser":                  "user",
	"getUser":                     "user",
	"getCurrentUser":              "current",
	"listUsers":                   "users",
	"getUserDeviceTypes":          "deviceTypes",
}
