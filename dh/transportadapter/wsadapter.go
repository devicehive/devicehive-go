package transportadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
)

type WSAdapter struct {
	transport transport.Transporter
}

type wsResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

func (a *WSAdapter) HandleResponseError(rawRes []byte) error {
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

func (a *WSAdapter) ResolveResource(resName string, data map[string]interface{}) (resource, method string) {
	if wsResources[resName] == "" {
		return resName, ""
	}

	return wsResources[resName], ""
}

func (a *WSAdapter) BuildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	return rawData
}

func (a *WSAdapter) ExtractResponsePayload(resourceName string, rawRes []byte) []byte {
	payloadKey := wsResponsePayloads[resourceName]
	if payloadKey == "" {
		return rawRes
	}

	res := make(map[string]json.RawMessage)
	json.Unmarshal(rawRes, &res)

	return res[payloadKey]
}

var wsResponsePayloads = map[string]string{
	"getConfig":          "configuration",
	"putConfig":          "configuration",
	"deleteConfig":       "configuration",
	"apiInfo":            "info",
	"apiInfoCluster":     "clusterInfo",
	"listCommands":       "commands",
	"insertCommand":      "command",
	"listNotifications":  "notifications",
	"insertNotification": "notification",
	"getDevice":          "device",
	"commandEvent":       "command",
	"notificationEvent":  "notification",
	"listDevices":        "devices",
}
