package dh

import (
	"github.com/devicehive/devicehive-go/internal/transport"
	"github.com/devicehive/devicehive-go/dh/transportadapter"
)

type Client struct {
	transport                 transport.Transporter
	transportAdapter          transportadapter.TransportAdapter
	accessToken               string
	refreshToken              string
	login                     string
	password                  string
	PollingWaitTimeoutSeconds int
}

func (c *Client) authenticate(token string) (result bool, err *Error) {
	if c.transport.IsHTTP() {
		c.accessToken = token
		return true, nil
	} else {
		_, err := c.request("auth", map[string]interface{}{
			"token": token,
		})

		if err != nil {
			return false, err
		}

		return true, nil
	}
}

func (c *Client) subscribe(resourceName string, params *SubscribeParams) (tspChan chan []byte, subscriptionId string, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	data, jsonErr := params.Map()
	if jsonErr != nil {
		return nil, "", &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	resource, tspReqParams := c.prepareRequestData(resourceName, data)
	if resource == "" {
		return nil, "", &Error{name: InvalidRequestErr, reason: "unknown resource name"}
	}

	tspReqParams.WaitTimeoutSeconds = c.PollingWaitTimeoutSeconds

	tspChan, subscriptionId, tspErr := c.transport.Subscribe(resource, tspReqParams)
	if tspErr != nil {
		return nil, "", newTransportErr(tspErr)
	}

	return tspChan, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	if c.transport.IsWS() {
		_, err := c.request(resourceName, map[string]interface{}{
			"subscriptionId": subscriptionId,
		})

		if err != nil {
			return err
		}
	}

	c.transport.Unsubscribe(subscriptionId)

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resource, tspReqParams := c.prepareRequestData(resourceName, data)

	if resource == "" {
		return nil, &Error{name: InvalidRequestErr, reason: "unknown resource name"}
	}

	resBytes, tspErr := c.transport.Request(resource, tspReqParams, Timeout)

	if tspErr != nil {
		return nil, newTransportErr(tspErr)
	}

	err = c.handleResponseError(resBytes)

	return resBytes, err
}

func (c *Client) handleResponseError(resBytes []byte) (err *Error) {
	rawErr := c.transportAdapter.HandleResponseError(resBytes)
	if rawErr != nil {
		return &Error{ServiceErr, rawErr.Error()}
	}

	return nil
}

func (c *Client) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, method := c.transportAdapter.ResolveResource(resourceName, data)
	reqData := c.transportAdapter.BuildRequestData(resourceName, data)
	reqParams = c.createRequestParams(method, reqData)

	return resource, reqParams
}

func (c *Client) createRequestParams(method string, reqData interface{}) *transport.RequestParams {
	tspReqParams := &transport.RequestParams{
		Data: reqData,
	}

	if c.transport.IsHTTP() {
		if method != "" {
			tspReqParams.Method = method
		}

		tspReqParams.AccessToken = c.accessToken
	}

	return tspReqParams
}
