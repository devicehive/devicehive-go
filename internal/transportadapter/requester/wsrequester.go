package requester

import (
	"github.com/devicehive/devicehive-go/internal/requestparams"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transportadapter/responsehandler"
	"time"
)

func NewWSRequester(tsp *transport.WS) *WSRequester {
	return &WSRequester{
		transport: tsp,
	}
}

type WSRequester struct {
	transport transport.Transporter
}

func (r *WSRequester) Request(resourceName string, data map[string]interface{}, timeout time.Duration, accessToken string) ([]byte, error) {
	resource, tspReqParams := r.PrepareRequestData(resourceName, data, accessToken)

	resBytes, tspErr := r.transport.Request(resource, tspReqParams, timeout)
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

func (r *WSRequester) PrepareRequestData(resourceName string, data map[string]interface{}, accessToken string) (resource string, reqParams *requestparams.RequestParams) {
	resource, _ = r.resolveResource(resourceName, data)
	reqParams = &requestparams.RequestParams{
		Data: data,
	}

	return resource, reqParams
}

func (r *WSRequester) resolveResource(resName string, data map[string]interface{}) (resource, method string) {
	if wsResources[resName] == "" {
		return resName, ""
	}

	return wsResources[resName], ""
}
