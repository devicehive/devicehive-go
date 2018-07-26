package requester

import (
	"github.com/devicehive/devicehive-go/internal/requester/httputils"
	"github.com/devicehive/devicehive-go/internal/resourcenames"
	"github.com/devicehive/devicehive-go/internal/responsehandler"
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/internal/transport/apirequests"
	"time"
)

func NewHTTPRequester(tsp transport.Transporter) *HTTPRequester {
	return &HTTPRequester{
		transport: tsp,
	}
}

type HTTPRequester struct {
	transport transport.Transporter
}

func (r *HTTPRequester) Request(resourceName string, data map[string]interface{}, timeout time.Duration, accessToken string) ([]byte, error) {
	resource, tspReqParams := r.PrepareRequestData(resourceName, data, accessToken)

	resBytes, tspErr := r.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err := responsehandler.HTTPHandleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	return resBytes, nil
}

func (r *HTTPRequester) PrepareRequestData(resourceName string, data map[string]interface{}, accessToken string) (resource string, reqParams *apirequests.RequestParams) {
	resource, method := r.ResolveResource(resourceName, data)
	reqData := r.buildRequestData(resourceName, data)
	reqParams = &apirequests.RequestParams{
		Data:   reqData,
		Method: method,
	}

	if resourceName != resourcenames.TokenRefresh && resourceName != resourcenames.TokenByCreds {
		reqParams.AccessToken = accessToken
	}

	return resource, reqParams
}

func (r *HTTPRequester) ResolveResource(resName string, data map[string]interface{}) (resource, method string) {
	rsrc, ok := httpResources[resName]

	if !ok {
		return resName, ""
	}

	queryParams := httputils.PrepareQueryParams(data)

	resource = httputils.PrepareHttpResource(rsrc[0], queryParams)
	method = rsrc[1]

	queryString := httputils.CreateQueryString(httpResourcesQueryParams, resName, queryParams)
	if queryString != "" {
		resource += "?" + queryString
	}

	return resource, method
}

func (r *HTTPRequester) buildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	payloadBuilder, ok := httpRequestPayloadBuilders[resourceName]

	if ok {
		return payloadBuilder(rawData)
	}

	return rawData
}
