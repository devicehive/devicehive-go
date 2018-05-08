package transport

import (
	"bytes"
	"encoding/json"
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
		client:        &http.Client{},
		url:           u,
		subscriptions: make(clientsMap),
	}, nil
}

type httpTsp struct {
	client        *http.Client
	url           *url.URL
	subscriptions clientsMap
}

func (t *httpTsp) IsHTTP() bool {
	return true
}

func (t *httpTsp) IsWS() bool {
	return false
}

func (t *httpTsp) Request(resource string, params *RequestParams, timeout time.Duration) (rawRes []byte, err *Error) {
	addr, err := t.createRequestAddr(resource)
	if err != nil {
		return nil, err
	}

	t.setTimeout(timeout)
	method := t.getRequestMethod(params)

	var req *http.Request
	var reqErr error
	if method != "GET" {
		reqDataReader, err := t.createRequestDataReader(params)
		if err != nil {
			return nil, err
		}

		req, reqErr = http.NewRequest(method, addr, reqDataReader)
	} else {
		req, reqErr = http.NewRequest(method, addr, nil)
	}

	if reqErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: reqErr.Error()}
	}

	if params != nil && params.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer "+params.AccessToken)
	}

	return t.doRequest(req)
}

func (t *httpTsp) setTimeout(timeout time.Duration) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	t.client.Timeout = timeout
}

func (t *httpTsp) getRequestMethod(params *RequestParams) string {
	if params == nil || params.Method == "" {
		return defaultHTTPMethod
	}

	return params.Method
}

func (t *httpTsp) createRequestDataReader(params *RequestParams) (dataReader *bytes.Reader, err *Error) {
	var rawReqData []byte

	if params != nil && params.Data != nil {
		var err error
		rawReqData, err = json.Marshal(params.Data)

		if err != nil {
			return nil, &Error{name: InvalidRequestErr, reason: err.Error()}
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
			return "", &Error{name: InvalidRequestErr, reason: urlErr.Error()}
		}
	}

	return u.String(), nil
}

func (t *httpTsp) doRequest(req *http.Request) (rawRes []byte, err *Error) {
	res, resErr := t.client.Do(req)

	if resErr != nil {
		if isTimeoutErr(resErr) {
			return nil, &Error{name: TimeoutErr, reason: resErr.Error()}
		}

		return nil, &Error{name: InvalidRequestErr, reason: resErr.Error()}
	}
	defer res.Body.Close()

	rawRes, rErr := ioutil.ReadAll(res.Body)

	if rErr != nil {
		return nil, &Error{name: InvalidResponseErr, reason: rErr.Error()}
	}

	return rawRes, nil
}

func (t *httpTsp) Subscribe(resource string, params *RequestParams) (eventChan chan []byte, subscriptionId string, err *Error) {
	subscriptionId = strconv.FormatInt(rand.Int63(), 10)

	subs := t.subscriptions.createSubscriber(subscriptionId)

	go func() {
		done := make(chan struct{})
		resChan := t.poll(resource, params, done)

		loop:
			for {
				select {
				case res := <-resChan:
					subs.data <- res
				case <-subs.signal:
					close(done)
					break loop
				}
			}
	}()

	return subs.data, subscriptionId, nil
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
	subscriber, ok := t.subscriptions.get(subscriptionId)

	if ok {
		subscriber.close()
		t.subscriptions.delete(subscriptionId)
	}
}
