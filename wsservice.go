package devicehive_go

import (
	"github.com/devicehive/devicehive-go/transport"
	"github.com/devicehive/devicehive-go/transportadapter"
)

func WSConnect(url string) (c *WSClient, err *Error) {
	tsp, tspErr := transport.Create(url)
	if tspErr != nil {
		return nil, &Error{name: ConnectionFailedErr, reason: tspErr.Error()}
	}

	if tsp.IsHTTP() {
		return nil, &Error{name: WrongURLErr, reason: "ws:// protocol is required"}
	}

	c = &WSClient{
		transportAdapter: transportadapter.New(tsp).(*transportadapter.WSAdapter),
		DataChan:         make(chan []byte),
		ErrorChan:        make(chan error),
	}

	return c, nil
}
