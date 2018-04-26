package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func newHTTP(addr string) (tsp *httpTsp, err error) {
	u, err := url.Parse(addr)

	if err != nil {
		return nil, err
	}

	return &httpTsp{
		client: &http.Client{},
		url:    u,
	}, nil
}

type httpTsp struct {
	client *http.Client
	url    *url.URL
}

func (t *httpTsp) Request(resource string, data devicehiveData, timeout time.Duration) (rawRes []byte, err *Error) {
	t.setTimeout(timeout)
	method := t.getRequestMethod(data)

	reqDataReader, err := t.createRequestDataReader(data)
	if err != nil {
		return nil, err
	}

	addr, err := t.createRequestAddr(data)
	if err != nil {
		return
	}

	req, reqErr := http.NewRequest(method, addr, reqDataReader)
	if reqErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: reqErr.Error()}
	}

	return t.doRequest(req)
}

func (t *httpTsp) setTimeout(timeout time.Duration) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	t.client.Timeout = timeout
}

func (t *httpTsp) getRequestMethod(data devicehiveData) string {
	defaultMethod := "GET"

	if data == nil {
		return defaultMethod
	}

	if _, ok := data["method"]; !ok {
		return defaultMethod
	}

	if m, ok := data["method"].(string); ok {
		return m
	}

	return defaultMethod
}

func (t *httpTsp) createRequestDataReader(data devicehiveData) (dataReader *bytes.Reader, err *Error) {
	var reqData []byte

	if reqData, ok := data["request"]; ok {
		var err error
		reqData, err = json.Marshal(reqData)

		if err != nil {
			return nil, &Error{name: InvalidRequestErr, reason: err.Error()}
		}
	} else {
		reqData = []byte("{}")
	}

	return bytes.NewReader(reqData), nil
}

func (t *httpTsp) createRequestAddr(data devicehiveData) (addr string, err *Error) {
	u := t.url

	if resource, ok := data["resource"].(string); ok {
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

func (t *httpTsp) Subscribe(subscriptionId string) (eventChan chan []byte) {
	return nil
}

func (t *httpTsp) Unsubscribe(subscriptionId string) {}
