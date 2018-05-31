package transport

import (
	"strings"
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

func Create(url string) (transport Transporter, err error) {
	if strings.Contains(url, "http") {
		return newHTTP(url)
	}

	return newWS(url)
}
