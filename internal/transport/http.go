package transport

import (
	"bytes"
	"encoding/json"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultHTTPMethod = "GET"
)

func newHTTP(addr string) (tsp *httpTsp, err error) {
	if addr[len(addr)-1:] != "/" {
		addr += "/"
	}

	u, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return &httpTsp{
		url:           u,
		subscriptions: apirequests.NewClientsMap(),
	}, nil
}

type httpTsp struct {
	url           *url.URL
	subscriptions *apirequests.PendingRequestsMap
}

func (t *httpTsp) IsHTTP() bool {
	return true
}

func (t *httpTsp) IsWS() bool {
	return false
}

func (t *httpTsp) Request(resource string, params *RequestParams, timeout time.Duration) (rawRes []byte, err *Error) {
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

func (t *httpTsp) getRequestMethod(params *RequestParams) string {
	if params == nil || params.Method == "" {
		return defaultHTTPMethod
	}

	return params.Method
}

func (t *httpTsp) createRequest(method, addr string, params *RequestParams) (req *http.Request, err error) {
	if method == "GET" {
		return http.NewRequest(method, addr, nil)
	}

	reqDataReader, err := t.createRequestDataReader(params)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, addr, reqDataReader)
}

func (t *httpTsp) createRequestDataReader(params *RequestParams) (dataReader *bytes.Reader, err error) {
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

func (t *httpTsp) createRequestAddr(resource string) (addr string, err *Error) {
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

func (t *httpTsp) addRequestHeaders(req *http.Request, params *RequestParams) {
	if params != nil && params.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+params.AccessToken)
	}
}

func (t *httpTsp) doRequest(client *http.Client, req *http.Request) (rawRes []byte, err *Error) {
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

func (t *httpTsp) Subscribe(resource string, params *RequestParams) (eventChan chan []byte, subscriptionId string, err *Error) {
	subscriptionId = strconv.FormatInt(rand.Int63(), 10)

	subs := t.subscriptions.CreateSubscription(subscriptionId)

	go func() {
		done := make(chan struct{})
		resChan := t.poll(resource, params, done)

	loop:
		for {
			select {
			case res := <-resChan:
				subs.Data <- res
			case <-subs.Signal:
				close(done)
				break loop
			}
		}
	}()

	return subs.Data, subscriptionId, nil
}

func (t *httpTsp) poll(resource string, params *RequestParams, done chan struct{}) (resChan chan []byte) {
	resChan = make(chan []byte)

	var timeout time.Duration
	if params == nil || params.WaitTimeoutSeconds == 0 {
		timeout = DefaultTimeout
	} else {
		timeout = time.Duration(params.WaitTimeoutSeconds) * time.Second * 2
	}

	go func() {
	loop:
		for {
			res, err := t.Request(resource, params, timeout)
			if err != nil {
				log.Printf("Subscription poll request failed for resource: %s, error: %s", resource, err)
				continue
			}

			select {
			case <-done:
				break loop
			case resChan <- res:
			}
		}
	}()

	return resChan
}

func (t *httpTsp) Unsubscribe(subscriptionId string) {
	subscriber, ok := t.subscriptions.Get(subscriptionId)

	if ok {
		subscriber.Close()
		t.subscriptions.Delete(subscriptionId)
	}
}
