// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transportadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
)

func newWSAdapter(tsp *transport.WS) *WSAdapter {
	a := &WSAdapter{
		transport: tsp,
	}

	tsp.AfterReconnection(func() {
		err := a.authenticatedResubscribe()

		if err != nil && err.Error() == "401 token expired" {
			tok, err := a.RefreshToken()
			if err != nil {
				tsp.TerminateRequests(err)
				return
			}

			a.accessToken = tok
			err = a.authenticatedResubscribe()
			if err != nil {
				tsp.TerminateRequests(err)
				return
			}
		} else if err != nil {
			tsp.TerminateRequests(err)
		}
	})

	return a
}

type WSAdapter struct {
	transport    *transport.WS
	accessToken  string
	login        string
	password     string
	refreshToken string
}

type wsResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

func (a *WSAdapter) SetCreds(login, password string) {
	a.login = login
	a.password = password
}

func (a *WSAdapter) SetRefreshToken(refTok string) {
	a.refreshToken = refTok
}

func (a *WSAdapter) authenticatedResubscribe() error {
	res, err := a.Authenticate(a.accessToken, 0)
	if res {
		a.transport.Resubscribe()
		return nil
	}

	return err
}

func (a *WSAdapter) Authenticate(token string, timeout time.Duration) (bool, error) {
	_, err := a.Request("auth", map[string]interface{}{
		"token": token,
	}, timeout)

	if err != nil {
		return false, err
	}

	a.accessToken = token
	return true, nil
}

func (a *WSAdapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, data)

	resBytes, tspErr := a.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err := a.handleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	resBytes = a.extractResponsePayload(resourceName, resBytes)

	return resBytes, nil
}

func (a *WSAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, params)

	tspSubs, subscriptionId, tspErr := a.transport.Subscribe(resource, tspReqParams)
	if tspErr != nil {
		return nil, "", tspErr
	}

	subscription = a.transformSubscription(resourceName, tspSubs)

	return subscription, subscriptionId, nil
}

func (a *WSAdapter) transformSubscription(resourceName string, subs *transport.Subscription) *transport.Subscription {
	dataChan := make(chan []byte)

	go func() {
		for d := range subs.DataChan {
			resErr := a.handleResponseError(d)
			if resErr != nil {
				subs.ErrChan <- resErr
			} else {
				data := a.extractResponsePayload(resourceName+"Event", d)
				dataChan <- data
			}
		}

		close(dataChan)
	}()

	transSubs := &transport.Subscription{
		DataChan: dataChan,
		ErrChan:  subs.ErrChan,
	}

	return transSubs
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

func (a *WSAdapter) RefreshToken() (accessToken string, err error) {
	if a.refreshToken == "" {
		accessToken, _, err = a.tokensByCreds(a.login, a.password)
		return accessToken, err
	}

	return a.accessTokenByRefresh(a.refreshToken)
}

func (a *WSAdapter) tokensByCreds(login, pass string) (accessToken, refreshToken string, err error) {
	rawRes, err := a.Request("tokenByCreds", map[string]interface{}{
		"login":    login,
		"password": pass,
	}, 0)

	if err != nil {
		return "", "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", "", parseErr
	}

	return tok.Access, tok.Refresh, nil
}

func (a *WSAdapter) accessTokenByRefresh(refreshToken string) (accessToken string, err error) {
	rawRes, err := a.Request("tokenRefresh", map[string]interface{}{
		"refreshToken": refreshToken,
	}, 0)

	if err != nil {
		return "", err
	}

	tok := &token{}
	parseErr := json.Unmarshal(rawRes, tok)

	if parseErr != nil {
		return "", parseErr
	}

	return tok.Access, nil
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

func (a *WSAdapter) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *apirequests.RequestParams) {
	resource, _ = a.resolveResource(resourceName, data)
	reqData := a.buildRequestData(resourceName, data)
	reqParams = &apirequests.RequestParams{
		Data: reqData,
	}

	return resource, reqParams
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
