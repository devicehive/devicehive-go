package dh

import (
	"encoding/json"
	"github.com/devicehive/devicehive-go/dh/transportadapter"
	"github.com/devicehive/devicehive-go/internal/transport"
)

type Client struct {
	transport                 transport.Transporter
	transportAdapter          transportadapter.TransportAdapter
	refreshToken              string
	login                     string
	password                  string
	PollingWaitTimeoutSeconds int
}

func (c *Client) authenticate(token string) (result bool, err *Error) {
	result, rawErr := c.transportAdapter.Authenticate(token, Timeout)

	if rawErr != nil {
		return false, newError(rawErr)
	}

	return true, nil
}

func (c *Client) subscribe(resourceName string, params *SubscribeParams) (tspChan chan []byte, subscriptionId string, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	data, jsonErr := params.Map()
	if jsonErr != nil {
		return nil, "", &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	tspChan, subscriptionId, rawErr := c.transportAdapter.Subscribe(resourceName, c.PollingWaitTimeoutSeconds, data)
	if rawErr != nil {
		return nil, "", newTransportErr(rawErr)
	}

	return tspChan, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	err := c.transportAdapter.Unsubscribe(resourceName, subscriptionId, Timeout)
	if err != nil {
		return newError(err)
	}

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resBytes, rawErr := c.transportAdapter.Request(resourceName, data, Timeout)

	if rawErr != nil && rawErr.Error() == "401 token expired" {
		resBytes, err = c.refreshRetry(resourceName, data)
		if err != nil {
			return nil, err
		}
	} else {
		err = newError(rawErr)
	}

	return resBytes, err
}

func (c *Client) refreshRetry(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	accessToken, err := c.RefreshToken()
	if err != nil {
		return nil, err
	}

	res, err := c.authenticate(accessToken)
	if !res || err != nil {
		return nil, err
	}

	resBytes, rawErr := c.transportAdapter.Request(resourceName, data, Timeout)
	if rawErr != nil {
		return nil, newError(rawErr)
	}

	return resBytes, nil
}

func (c *Client) getModel(resourceName string, model interface{}, data map[string]interface{}) *Error {
	rawRes, err := c.request(resourceName, data)

	if err != nil {
		return err
	}

	parseErr := json.Unmarshal(rawRes, model)
	if parseErr != nil {
		return newJSONErr()
	}

	return nil
}
