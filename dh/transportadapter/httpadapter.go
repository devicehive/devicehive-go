package transportadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
	"time"
)

type HTTPAdapter struct {
	transport transport.Transporter
}

type httpResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (a *HTTPAdapter) HandleResponseError(rawRes []byte) error {
	if len(rawRes) == 0 {
		return nil
	}

	if isJSONArray(rawRes) {
		return nil
	}

	httpRes, err := a.formatHTTPResponse(rawRes)
	if httpRes == nil && err == nil {
		return nil
	} else if err != nil {
		return err
	}

	if httpRes.Status >= 400 {
		errMsg := strings.ToLower(httpRes.Message)
		errCode := httpRes.Status
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return errors.New(r)
	}

	return nil
}

func (a *HTTPAdapter) formatHTTPResponse(rawRes []byte) (httpRes *httpResponse, err error) {
	res := make(map[string]interface{})
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, err
	}

	if _, ok := res["message"]; !ok {
		return nil, nil
	}

	httpRes = &httpResponse{
		Message: res["message"].(string),
	}
	if e, ok := res["error"].(float64); ok {
		httpRes.Status = int(e)
	} else {
		httpRes.Status = int(res["status"].(float64))
	}

	return httpRes, nil
}

func (a *HTTPAdapter) ResolveResource(resName string, data map[string]interface{}) (resource, method string) {
	rsrc, ok := httpResources[resName]

	if !ok {
		return resName, ""
	}

	queryParams := prepareQueryParams(data)

	resource = prepareHttpResource(rsrc[0], queryParams)
	method = rsrc[1]

	queryString := createQueryString(httpResourcesQueryParams, resName, queryParams)
	if queryString != "" {
		resource += "?" + queryString
	}

	return resource, method
}

func (a *HTTPAdapter) BuildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	payloadBuilder, ok := httpRequestPayloadBuilders[resourceName]

	if ok {
		return payloadBuilder(rawData)
	}

	return rawData
}

func (a *HTTPAdapter) ExtractResponsePayload(resourceName string, rawRes []byte) []byte {
	return rawRes
}

func (a *HTTPAdapter) Request(resourceName, accessToken string, data map[string]interface{}, timeout time.Duration) (res []byte, err error) {
	resource, tspReqParams := a.prepareRequestData(resourceName, accessToken, data)

	resBytes, tspErr := a.transport.Request(resource, tspReqParams, timeout)
	if tspErr != nil {
		return nil, tspErr
	}

	err = a.HandleResponseError(resBytes)
	if err != nil {
		return nil, err
	}

	resBytes = a.ExtractResponsePayload(resourceName, resBytes)

	return resBytes, nil
}

func (a *HTTPAdapter) prepareRequestData(resourceName, accessToken string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, method := a.ResolveResource(resourceName, data)
	reqData := a.BuildRequestData(resourceName, data)
	reqParams = &transport.RequestParams{
		Data: reqData,
		Method: method,
		AccessToken: accessToken,
	}

	return resource, reqParams
}
