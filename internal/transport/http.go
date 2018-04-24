package transport

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func newHTTP(addr string) *httpTsp {
	return &httpTsp{
		client: &http.Client{},
		addr:   addr,
	}
}

type httpTsp struct {
	client *http.Client
	addr   string
}

func (t *httpTsp) Request(data devicehiveData, timeout time.Duration) (rawRes []byte, err *Error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

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
	req, reqErr := http.NewRequest(data["method"].(string), t.addr, dataReader)

	if reqErr != nil {
		return nil, &Error{name: InvalidRequestErr, reason: reqErr.Error()}
	}

	res, resErr := t.client.Do(req)

	if resErr != nil {
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
