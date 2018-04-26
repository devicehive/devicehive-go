package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"net/url"
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
	url   *url.URL
}

func (t *httpTsp) Request(resource string, data devicehiveData, timeout time.Duration) (rawRes []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	t.client.Timeout = timeout

	if _, ok := data["method"]; !ok {
		data["method"] = "GET"
	}

	var reqRawData []byte
	if reqData, ok := data["request"]; ok {
		var err error
		reqRawData, err = json.Marshal(reqData)

		if err != nil {
			return nil, &Error{name: InvalidRequestErr, reason: err.Error()}
		}
	} else {
		reqRawData = []byte("{}")
	}

	dataReader := bytes.NewReader(reqRawData)

	addr := t.url

	if resource, ok := data["resource"].(string); ok {
		var urlErr error
		addr, urlErr = t.url.Parse(resource)

		if urlErr != nil {
			return nil, &Error{name: InvalidRequestErr, reason: urlErr.Error()}
		}
	}

	req, reqErr := http.NewRequest(data["method"].(string), addr.String(), dataReader)

	if reqErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: reqErr.Error()}
	}

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
