// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transportadapter

import (
	"time"

	"github.com/devicehive/devicehive-go/internal/authmanager"
	"github.com/devicehive/devicehive-go/internal/requester"
	"github.com/devicehive/devicehive-go/internal/responsehandler"
	"github.com/devicehive/devicehive-go/internal/transport"
)

func newWSAdapter(tsp *transport.WS) *WSAdapter {
	reqstr := requester.NewWSRequester(tsp)
	a := &WSAdapter{
		transport: tsp,
		reqstr:    reqstr,
		authMng:   authmanager.New(reqstr),
	}

	tsp.AfterReconnection(func() {
		err := a.authenticatedResubscribe()

		if err != nil && err.Error() == TokenExpiredErr {
			tok, err := a.RefreshToken()
			if err != nil {
				tsp.TerminateRequests(err)
				return
			}

			a.authMng.SetAccessToken(tok)
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
	transport *transport.WS
	authMng   *authmanager.AuthManager
	reqstr    *requester.WSRequester
}

func (a *WSAdapter) SetCreds(login, password string) {
	a.authMng.SetCreds(login, password)
}

func (a *WSAdapter) SetRefreshToken(refTok string) {
	a.authMng.SetRefreshToken(refTok)
}

func (a *WSAdapter) Authenticate(token string, timeout time.Duration) (bool, error) {
	_, err := a.Request("auth", map[string]interface{}{
		"token": token,
	}, timeout)

	if err != nil {
		return false, err
	}

	a.authMng.SetAccessToken(token)
	return true, nil
}

func (a *WSAdapter) authenticatedResubscribe() error {
	res, err := a.Authenticate(a.authMng.AccessToken(), 0)
	if res {
		a.transport.Resubscribe()
		return nil
	}

	return err
}

func (a *WSAdapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	return a.reqstr.Request(resourceName, data, timeout, a.authMng.AccessToken())
}

func (a *WSAdapter) Subscribe(resourceName string, pollingWaitTimeoutSeconds int, params map[string]interface{}) (subscription *transport.Subscription, subscriptionId string, err *transport.Error) {
	resource, tspReqParams := a.reqstr.PrepareRequestData(resourceName, params, a.authMng.AccessToken())

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
	return a.authMng.RefreshToken()
}
