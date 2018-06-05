package transportadapter

import "github.com/devicehive/devicehive-go/transport"

type RequestResponseHandler interface {
	handleResponseError(rawRes []byte) error
	extractResponsePayload(resourceName string, rawRes []byte) []byte
	prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams)
}
