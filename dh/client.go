package dh

import (
	"github.com/devicehive/devicehive-go/dh/transportadapter"
	"github.com/devicehive/devicehive-go/internal/transport"
	"encoding/json"
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

	tspChan, subscriptionId, rawErr := c.transportAdapter.Subscribe(resourceName, c.accessToken, c.PollingWaitTimeoutSeconds, data)
	if rawErr != nil {
		return nil, "", newTransportErr(rawErr)
	}

	return tspChan, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	err := c.transportAdapter.Unsubscribe(resourceName, c.accessToken, subscriptionId, Timeout)

	if err != nil {
		switch err.(type) {
		case *transport.Error:
			return newTransportErr(err.(*transport.Error))
		default:
			return &Error{ServiceErr, err.Error()}
		}
	}

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resBytes, rawErr := c.transportAdapter.Request(resourceName, c.accessToken, data, Timeout)

	if rawErr != nil {
		switch rawErr.(type) {
		case *transport.Error:
			return nil, newTransportErr(rawErr.(*transport.Error))
		default:
			return nil, &Error{ServiceErr, rawErr.Error()}
		}
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
