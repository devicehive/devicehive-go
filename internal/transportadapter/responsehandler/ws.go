package responsehandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/resourcenames"
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
	resourcenames.GetConfig:          "configuration",
	resourcenames.PutConfig:          "configuration",
	resourcenames.DeleteConfig:       "configuration",
	resourcenames.ApiInfo:            "info",
	resourcenames.ClusterInfo:        "clusterInfo",
	resourcenames.ListCommands:       "commands",
	resourcenames.InsertCommand:      "command",
	resourcenames.ListNotifications:  "notifications",
	resourcenames.InsertNotification: "notification",
	"subscribeNotificationsEvent":    "notification",
	"subscribeCommandsEvent":         "command",
	resourcenames.GetDevice:          "device",
	"commandEvent":                   "command",
	"notificationEvent":              "notification",
	resourcenames.ListDevices:        "devices",
	resourcenames.InsertNetwork:      "network",
	resourcenames.GetNetwork:         "network",
	resourcenames.ListNetworks:       "networks",
	resourcenames.InsertDeviceType:   "deviceType",
	resourcenames.GetDeviceType:      "deviceType",
	resourcenames.ListDeviceTypes:    "deviceTypes",
	resourcenames.CreateUser:         "user",
	resourcenames.GetUser:            "user",
	resourcenames.GetCurrentUser:     "current",
	resourcenames.ListUsers:          "users",
	resourcenames.GetUserDeviceTypes: "deviceTypes",
}
