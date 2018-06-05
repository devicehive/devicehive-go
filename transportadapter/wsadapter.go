package transportadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devicehive/devicehive-go/transport"
	"strings"
	"time"
)

type WSAdapter struct {
	*Adapter
}

type wsResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

func (a *WSAdapter) Authenticate(token string, timeout time.Duration) (result bool, err error) {
	_, err = a.Request("auth", map[string]interface{}{
		"token": token,
	}, timeout)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (a *WSAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (tspChan chan []byte, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, params)

	tspChan, subscriptionId, tspErr := a.transport.Subscribe(resource, tspReqParams)
	if tspErr != nil {
		return nil, "", tspErr
	}

	return tspChan, subscriptionId, nil
}

func (a *WSAdapter) Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error {
	_, err := a.Request(resourceName, map[string]interface{}{
		"subscriptionId": subscriptionId,
	}, timeout)

	if err != nil {
		return err
	}

	a.transport.Unsubscribe(subscriptionId)

	return nil
}

func (a *WSAdapter) handleResponseError(rawRes []byte) error {
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

func (a *WSAdapter) resolveResource(resName string, data map[string]interface{}) (resource, method string) {
	if wsResources[resName] == "" {
		return resName, ""
	}

	return wsResources[resName], ""
}

func (a *WSAdapter) buildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	return rawData
}

func (a *WSAdapter) extractResponsePayload(resourceName string, rawRes []byte) []byte {
	payloadKey := wsResponsePayloads[resourceName]
	if payloadKey == "" {
		return rawRes
	}

	res := make(map[string]json.RawMessage)
	json.Unmarshal(rawRes, &res)

	return res[payloadKey]
}

func (a *WSAdapter) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, _ = a.resolveResource(resourceName, data)
	reqData := a.buildRequestData(resourceName, data)
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
