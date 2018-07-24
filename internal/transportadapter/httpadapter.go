// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transportadapter

import (
	"encoding/json"
	"time"

	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter/requester"
	"github.com/devicehive/devicehive-go/internal/transportadapter/responsehandler"
)

type Timestamp struct {
	Value string `json:"timestamp"`
}

func newHTTPAdapter(tsp *transport.HTTP) *HTTPAdapter {
	return &HTTPAdapter{
		transport: tsp,
		reqstr: requester.NewHTTPRequester(tsp),
	}
}

type HTTPAdapter struct {
	transport   *transport.HTTP
	accessToken string
	refreshToken string
	login string
	password string
	reqstr       *requester.HTTPRequester
}

func (a *HTTPAdapter) SetCreds(login, password string) {
	a.login = login
	a.password = password
}

func (a *HTTPAdapter) SetRefreshToken(refTok string) {
	a.refreshToken = refTok
}

func (a *HTTPAdapter) Authenticate(token string, timeout time.Duration) (bool, error) {
	a.transport.SetPollingToken(token)
	a.accessToken = token
	return true, nil
}

func (a *HTTPAdapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	res, err := a.request(resourceName, data, timeout)
	if err != nil && err.Error() == TokenExpiredErr {
		res, err := a.refreshRetry(resourceName, data, timeout)
		return res, err
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *HTTPAdapter) refreshRetry(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	accessToken, err := a.RefreshToken()
	if err != nil {
		return nil, err
	}

	res, err := a.Authenticate(accessToken, 0)
	if !res || err != nil {
		return nil, err
	}

	return a.request(resourceName, data, 0)
}

func (a *HTTPAdapter) request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	return a.reqstr.Request(resourceName, data, timeout, a.accessToken)
}

func (a *HTTPAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.reqstr.PrepareRequestData(resourceName, params, a.accessToken)

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
		if resErr := responsehandler.HTTPHandleResponseError(data); resErr != nil {
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

	resource, _ := a.reqstr.ResolveResource(resourceName, params)

	a.transport.SetPollingResource(subscriptionId, resource)
}

func (a *HTTPAdapter) Unsubscribe(resourceName, subscriptionId string, timeout time.Duration) error {
	a.transport.Unsubscribe(subscriptionId)
	return nil
}

func (a *HTTPAdapter) RefreshToken() (accessToken string, err error) {
	if a.refreshToken == "" {
		accessToken, _, err = a.TokensByCreds(a.login, a.password)
		return accessToken, err
	}

	return a.AccessTokenByRefresh(a.refreshToken)
}

func (a *HTTPAdapter) TokensByCreds(login, pass string) (accessToken, refreshToken string, err error) {
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

func (a *HTTPAdapter) AccessTokenByRefresh(refreshToken string) (accessToken string, err error) {
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
