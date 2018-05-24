package transportadapter

import (
	"time"
	"github.com/devicehive/devicehive-go/internal/transport"
)

type Adapter struct {
	transport transport.Transporter
	adapter RequestResponseHandler
}

func (a *Adapter) Request(resourceName string, data map[string]interface{}, timeout time.Duration) (res []byte, err error) {
	resource, tspReqParams := a.adapter.prepareRequestData(resourceName, data)

	resBytes, tspErr := a.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err = a.adapter.handleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	resBytes = a.adapter.extractResponsePayload(resourceName, resBytes)

	return resBytes, nil
}
