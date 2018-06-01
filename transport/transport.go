package transport

import (
	"net/url"
	"time"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Transporter interface {
	Request(resource string, params *RequestParams, timeout time.Duration) (res []byte, err *Error)
	Subscribe(resource string, params *RequestParams) (eventChan chan []byte, subscriptionId string, err *Error)
	Unsubscribe(subscriptionId string)
	IsHTTP() bool
	IsWS() bool
}

func Create(addr string) (transport Transporter, err error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "http" || u.Scheme == "https" {
		return newHTTP(addr)
	}

	return newWS(addr)
}
