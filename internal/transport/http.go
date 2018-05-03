package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultHTTPMethod = "GET"
)

func newHTTP(addr string) (tsp *httpTsp, err error) {
	if addr[len(addr) - 1:] != "/" {
		addr += "/"
	}

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

func (t *httpTsp) IsHTTP() bool {
	return true
}

func (t *httpTsp) IsWS() bool {
	return false
}

func (t *httpTsp) Request(resource string, params *RequestParams, timeout time.Duration) (rawRes []byte, err *Error) {
	t.setTimeout(timeout)
	method := t.getRequestMethod(params)

	reqDataReader, err := t.createRequestDataReader(params)
	if err != nil {
		return nil, err
	}

	addr, err := t.createRequestAddr(resource)
	if err != nil {
		return
	}

	req, reqErr := http.NewRequest(method, addr, reqDataReader)
	if reqErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: reqErr.Error()}
	}

	if params != nil && params.AccessToken != "" {
		req.Header.Add("Authorization", "Bearer " + params.AccessToken)
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

func (t *httpTsp) Subscribe(subscriptionId string) (eventChan chan []byte) {
	return nil
}

func (t *httpTsp) Unsubscribe(subscriptionId string) {}
