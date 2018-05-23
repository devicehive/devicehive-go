package transportadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
	"time"
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

func (a *WSAdapter) Request(resourceName, accessToken string, data map[string]interface{}, timeout time.Duration) (res []byte, err error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, data)

	resBytes, tspErr := a.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err = a.HandleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	resBytes = a.ExtractResponsePayload(resourceName, resBytes)

	return resBytes, nil
}

func (a *WSAdapter) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, _ = a.ResolveResource(resourceName, data)
	reqData := a.BuildRequestData(resourceName, data)
	reqParams = &transport.RequestParams{
		Data: reqData,
	}

	return resource, reqParams
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
	"insertNetwork":      "network",
	"getNetwork":         "network",
	"listNetworks":       "networks",
	"insertDeviceType":   "deviceType",
	"getDeviceType":      "deviceType",
	"listDeviceTypes":    "deviceTypes",
	"createUser":         "user",
	"getUser":            "user",
	"getCurrentUser":     "current",
	"listUsers":          "users",
	"getUserDeviceTypes": "deviceTypes",
}
