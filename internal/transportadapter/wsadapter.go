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

func newWSAdapter(tsp *transport.WS) *WSAdapter {
	a := &WSAdapter{
		transport: tsp,
		reqstr: requester.New(tsp),
	}

	tsp.AfterReconnection(func() {
		err := a.authenticatedResubscribe()

		if err != nil && err.Error() == TokenExpiredErr {
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
	reqstr       requester.Requester
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
	return a.reqstr.Request(resourceName, data, timeout)
}

func (a *WSAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.reqstr.PrepareRequestData(resourceName, params)

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
			resErr := responsehandler.WSHandleResponseError(d)
			if resErr != nil {
				subs.ErrChan <- resErr
			} else {
				data := responsehandler.WSExtractResponsePayload(resourceName+"Event", d)
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
		accessToken, _, err = a.TokensByCreds(a.login, a.password)
		return accessToken, err
	}

	return a.AccessTokenByRefresh(a.refreshToken)
}

func (a *WSAdapter) TokensByCreds(login, pass string) (accessToken, refreshToken string, err error) {
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

func (a *WSAdapter) AccessTokenByRefresh(refreshToken string) (accessToken string, err error) {
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
