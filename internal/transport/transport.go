package transport

import (
	"strings"
	"time"
)

const (
	DefaultTimeout = 3 * time.Second
)

type Transporter interface {
	Request(data devicehiveData, timeout time.Duration) (res []byte, err *Error)
	Subscribe(subscriptionId string) (eventChan chan []byte)
	Unsubscribe(subscriptionId string)
}

func Create(url string) (transport Transporter, err error) {
	if strings.Contains(url, "http") {
		return newHTTP(url)
	}

	return newWS(url)
}
