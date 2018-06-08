// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/devicehive/devicehive-go/transport/apirequests"
)

const (
	defaultHTTPMethod = "GET"
)

func newHTTP(addr string) (*HTTP, error) {
	if addr[len(addr)-1:] != "/" {
		addr += "/"
	}

	u, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return &HTTP{
		url:           u,
		subscriptions: apirequests.NewClientsMap(),
	}, nil
}

type HTTP struct {
	url                *url.URL
	subscriptions      *apirequests.PendingRequestsMap
	pollingAccessToken string
}

func (t *HTTP) IsHTTP() bool {
	return true
}

func (t *HTTP) IsWS() bool {
	return false
}

func (t *HTTP) SetPollingToken(accessToken string) {
	t.pollingAccessToken = accessToken
}

func (t *HTTP) Request(resource string, params *RequestParams, timeout time.Duration) ([]byte, *Error) {
	client := &http.Client{}
	addr, err := t.createRequestAddr(resource)
	if err != nil {
		return nil, err
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}
	client.Timeout = timeout

	method := t.getRequestMethod(params)
	req, reqErr := t.createRequest(method, addr, params)
	if reqErr != nil {
		return nil, NewError(InvalidRequestErr, reqErr.Error())
	}

	t.addRequestHeaders(req, params)

	return t.doRequest(client, req)
}

func (t *HTTP) getRequestMethod(params *RequestParams) string {
	if params == nil || params.Method == "" {
		return defaultHTTPMethod
	}

	return params.Method
}

func (t *HTTP) createRequest(method, addr string, params *RequestParams) (*http.Request, error) {
	if method == "GET" {
		return http.NewRequest(method, addr, nil)
	}

	reqDataReader, err := t.createRequestDataReader(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, addr, reqDataReader)
}

func (t *HTTP) createRequestDataReader(params *RequestParams) (*bytes.Reader, error) {
	var rawReqData []byte

	if params != nil && params.Data != nil {
		var err error
		rawReqData, err = json.Marshal(params.Data)

		if err != nil {
			return nil, err
		}
	} else {
		rawReqData = []byte("{}")
	}

	return bytes.NewReader(rawReqData), nil
}

func (t *HTTP) createRequestAddr(resource string) (addr string, err *Error) {
	u := t.url

	if resource != "" {
		var urlErr error
		u, urlErr = t.url.Parse(resource)

		if urlErr != nil {
			return "", NewError(InvalidRequestErr, urlErr.Error())
		}
	}

	return u.String(), nil
}

func (t *HTTP) addRequestHeaders(req *http.Request, params *RequestParams) {
	if params != nil && params.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+params.AccessToken)
	}
}

func (t *HTTP) doRequest(client *http.Client, req *http.Request) ([]byte, *Error) {
	res, resErr := client.Do(req)

	if resErr != nil {
		if isTimeoutErr(resErr) {
			return nil, NewError(TimeoutErr, resErr.Error())
		}

		return nil, NewError(InvalidRequestErr, resErr.Error())
	}
	defer res.Body.Close()

	rawRes, rErr := ioutil.ReadAll(res.Body)

	if rErr != nil {
		return nil, NewError(InvalidResponseErr, rErr.Error())
	}

	return rawRes, nil
}

func (t *HTTP) Subscribe(resource string, params *RequestParams) (subscription *Subscription, subscriptionId string, err *Error) {
	subscriptionId = strconv.FormatInt(rand.Int63(), 10)

	subs := t.subscriptions.CreateRequest(subscriptionId)

	go func() {
		done := make(chan struct{})
		resChan, errChan := t.poll(resource, params, done)

	loop:
		for {
			select {
			case err := <-errChan:
				subs.Err <- err
			case res := <-resChan:
				subs.Data <- res
			case <-subs.Signal:
				close(subs.Data)
				close(subs.Err)
				close(done)
				break loop
			}
		}
	}()

	subscription = &Subscription{
		DataChan: subs.Data,
		ErrChan:  subs.Err,
	}

	return subscription, subscriptionId, nil
}

func (t *HTTP) poll(resource string, params *RequestParams, done chan struct{}) (resChan chan []byte, errChan chan error) {
	resChan = make(chan []byte)
	errChan = make(chan error)

	var timeout time.Duration
	if params == nil || params.WaitTimeoutSeconds == 0 {
		timeout = DefaultTimeout
	} else {
		timeout = time.Duration(params.WaitTimeoutSeconds) * time.Second * 2
	}

	go func() {
	loop:
		for {
			params.AccessToken = t.pollingAccessToken
			res, err := t.Request(resource, params, timeout)
			if err != nil {
				errChan <- err
				continue
			}

			select {
			case <-done:
				break loop
			case resChan <- res:
			}
		}
	}()

	return resChan, errChan
}

func (t *HTTP) Unsubscribe(subscriptionId string) {
	subscriber, ok := t.subscriptions.Get(subscriptionId)

	if ok {
		subscriber.Close()
		t.subscriptions.Delete(subscriptionId)
	}
}
