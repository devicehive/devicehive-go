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
)

type Timestamp struct {
	Value string `json:"timestamp"`
}

type HTTPAdapter struct {
	transport   *transport.HTTP
	accessToken string
}

type httpResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (a *HTTPAdapter) Authenticate(token string, timeout time.Duration) (bool, error) {
	a.transport.SetPollingToken(token)
	a.accessToken = token
	return true, nil
}

func (a *HTTPAdapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
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

func (a *HTTPAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, params)

	tspReqParams.WaitTimeoutSeconds = pollingWaitTimeoutSeconds

	tspSubs, subscriptionId, tspErr := a.transport.Subscribe(resource, tspReqParams)
	if tspErr != nil {
		return nil, "", tspErr
	}

	subscription = a.transformSubscription(resourceName, subscriptionId, params, tspSubs)

	return subscription, subscriptionId, nil
}

func (a *HTTPAdapter) transformSubscription(resourceName, subscriptionId string, params map[string]interface{}, subs *transport.Subscription) *transport.Subscription {
	dataChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
	loop:
		for {
			select {
			case d, ok := <-subs.DataChan:
				if !ok {
					break loop
				}

				list, err := a.handleSubscriptionEventData(d)
				if err != nil {
					errChan <- err
					continue
				}

				a.setResourceWithLastEntityTimestamp(resourceName, subscriptionId, params, list)
				subs.ContinuePolling()

				for _, data := range list {
					dataChan <- data
				}
			case err, ok := <-subs.ErrChan:
				if !ok {
					break loop
				}

				errChan <- err
				subs.ContinuePolling()
			}
		}

		close(dataChan)
		close(errChan)
	}()

	transSubs := &transport.Subscription{
		DataChan: dataChan,
		ErrChan:  errChan,
	}

	return transSubs
}

func (a *HTTPAdapter) handleSubscriptionEventData(data []byte) ([]json.RawMessage, error) {
	var list []json.RawMessage
	if err := json.Unmarshal(data, &list); err != nil {
		if resErr := a.handleResponseError(data); resErr != nil {
			return nil, resErr
		} else {
			return nil, err
		}
	}

	return list, nil
}

func (a *HTTPAdapter) setResourceWithLastEntityTimestamp(resourceName, subscriptionId string, params map[string]interface{}, list []json.RawMessage) {
	l := len(list)
	if l == 0 {
		return
	}

	timestamp := &Timestamp{}
	json.Unmarshal(list[l-1], timestamp)

	if timestamp.Value == "" {
		return
	}

	if params == nil {
		params = make(map[string]interface{})
	}
	params["timestamp"] = timestamp.Value

	resource, _ := a.resolveResource(resourceName, params)

	a.transport.SetPollingResource(subscriptionId, resource)
}

func (a *HTTPAdapter) Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error {
	a.transport.Unsubscribe(subscriptionId)
	return nil
}

func (a *HTTPAdapter) handleResponseError(rawRes []byte) error {
	if len(rawRes) == 0 {
		return nil
	}

	if isJSONArray(rawRes) {
		return nil
	}

	httpRes, err := a.formatHTTPResponse(rawRes)
	if httpRes == nil && err == nil {
		return nil
	} else if err != nil {
		return err
	}

	if httpRes.Status >= 400 {
		errMsg := strings.ToLower(httpRes.Message)
		errCode := httpRes.Status
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return errors.New(r)
	}

	return nil
}

func (a *HTTPAdapter) formatHTTPResponse(rawRes []byte) (*httpResponse, error) {
	res := make(map[string]interface{})
	err := json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, err
	}

	if _, ok := res["message"]; !ok {
		return nil, nil
	}

	httpRes := &httpResponse{
		Message: res["message"].(string),
	}
	if e, ok := res["error"].(float64); ok {
		httpRes.Status = int(e)
	} else {
		httpRes.Status = int(res["status"].(float64))
	}

	return httpRes, nil
}

func (a *HTTPAdapter) resolveResource(resName string, data map[string]interface{}) (resource, method string) {
	rsrc, ok := httpResources[resName]

	if !ok {
		return resName, ""
	}

	queryParams := prepareQueryParams(data)

	resource = prepareHttpResource(rsrc[0], queryParams)
	method = rsrc[1]

	queryString := createQueryString(httpResourcesQueryParams, resName, queryParams)
	if queryString != "" {
		resource += "?" + queryString
	}

	return resource, method
}

func (a *HTTPAdapter) buildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	payloadBuilder, ok := httpRequestPayloadBuilders[resourceName]

	if ok {
		return payloadBuilder(rawData)
	}

	return rawData
}

func (a *HTTPAdapter) extractResponsePayload(resourceName string, rawRes []byte) []byte {
	return rawRes
}

func (a *HTTPAdapter) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, method := a.resolveResource(resourceName, data)
	reqData := a.buildRequestData(resourceName, data)
	reqParams = &transport.RequestParams{
		Data:   reqData,
		Method: method,
	}

	if resourceName != "tokenRefresh" && resourceName != "tokenByCreds" {
		reqParams.AccessToken = a.accessToken
	}

	return resource, reqParams
}
