// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transportadapter

import (
	"encoding/json"
	"time"

	"github.com/devicehive/devicehive-go/internal/authmanager"
	"github.com/devicehive/devicehive-go/internal/requester"
	"github.com/devicehive/devicehive-go/internal/responsehandler"
	"github.com/devicehive/devicehive-go/internal/transport"
)

type Timestamp struct {
	Value string `json:"timestamp"`
}

func newHTTPAdapter(tsp *transport.HTTP) *HTTPAdapter {
	reqstr := requester.NewHTTPRequester(tsp)
	return &HTTPAdapter{
		transport: tsp,
		reqstr:    reqstr,
		authMng:   authmanager.New(reqstr),
	}
}

type HTTPAdapter struct {
	transport *transport.HTTP
	authMng   *authmanager.AuthManager
	reqstr    *requester.HTTPRequester
}

func (a *HTTPAdapter) SetCreds(login, password string) {
	a.authMng.SetCreds(login, password)
}

func (a *HTTPAdapter) SetRefreshToken(refTok string) {
	a.authMng.SetRefreshToken(refTok)
}

func (a *HTTPAdapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	res, err := a.request(resourceName, data, timeout)
	if a.isReauthNeeded(err) {
		res, err := a.refreshRetry(resourceName, data, timeout)
		return res, err
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *HTTPAdapter) refreshRetry(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	err := a.reauthenticate()
	if err != nil {
		return nil, err
	}

	return a.request(resourceName, data, 0)
}

func (a *HTTPAdapter) request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	return a.reqstr.Request(resourceName, data, timeout, a.authMng.AccessToken())
}

func (a *HTTPAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.reqstr.PrepareRequestData(resourceName, params, a.authMng.AccessToken())

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
				if a.isReauthNeeded(err) {
					err := a.reauthenticate()
					if err != nil {
						errChan <- err
					}
					break
				} else if err != nil {
					errChan <- err
					break
				}

				a.setResourceWithLastEntityTimestamp(resourceName, subscriptionId, params, list)

				for _, data := range list {
					dataChan <- data
				}
			case err, ok := <-subs.ErrChan:
				if !ok {
					break loop
				}

				errChan <- err
			}

			subs.ContinuePolling()
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

func (a *HTTPAdapter) reauthenticate() error {
	defer a.authMng.Reauth.Checkpoint()

	accessToken, err := a.RefreshToken()
	if err != nil {
		return err
	}

	_, err = a.Authenticate(accessToken, 0)
	return err
}

func (a *HTTPAdapter) RefreshToken() (accessToken string, err error) {
	return a.authMng.RefreshToken()
}

func (a *HTTPAdapter) Authenticate(token string, timeout time.Duration) (bool, error) {
	a.transport.SetPollingToken(token)
	a.authMng.SetAccessToken(token)
	return true, nil
}

func (a *HTTPAdapter) isReauthNeeded(err error) bool {
	return err != nil && err.Error() == TokenExpiredErr && a.authMng.Reauth.Needed()
}
