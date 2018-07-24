package requester

import (
	"time"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"github.com/devicehive/devicehive-go/internal/transportadapter/responsehandler"
)

type WSRequester struct {
	transport transport.Transporter
}

func (a *WSRequester) Request(resourceName string, data map[string]interface{}, timeout time.Duration) ([]byte, error) {
	resource, tspReqParams := a.PrepareRequestData(resourceName, data)

	resBytes, tspErr := a.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err := responsehandler.WSHandleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	resBytes = responsehandler.WSExtractResponsePayload(resourceName, resBytes)

	return resBytes, nil
}

func (a *WSRequester) PrepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *apirequests.RequestParams) {
	resource, _ = a.resolveResource(resourceName, data)
	reqParams = &apirequests.RequestParams{
		Data: data,
	}

	return resource, reqParams
}

func (a *WSRequester) resolveResource(resName string, data map[string]interface{}) (resource, method string) {
	if wsResources[resName] == "" {
		return resName, ""
	}

	return wsResources[resName], ""
}
