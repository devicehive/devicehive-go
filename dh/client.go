package dh

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
)

type Client struct {
	tsp          transport.Transporter
	accessToken  string
	refreshToken string
	login        string
	password     string
	PollingWaitTimeoutSeconds int
}

func (c *Client) authenticate(token string) (result bool, err *Error) {
	if c.tsp.IsHTTP() {
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

	tspChan, subscriptionId, tspErr := c.tsp.Subscribe(resource, tspReqParams)

	if tspErr != nil {
		return nil, "", newTransportErr(tspErr)
	}

	return tspChan, subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	if c.tsp.IsWS() {
		_, err := c.request(resourceName, map[string]interface{}{
			"subscriptionId": subscriptionId,
		})

		if err != nil {
			return err
		}
	}

	c.tsp.Unsubscribe(subscriptionId)

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resource, tspReqParams := c.prepareRequestData(resourceName, data)

	if resource == "" {
		return nil, &Error{name: InvalidRequestErr, reason: "unknown resource name"}
	}

	resBytes, tspErr := c.tsp.Request(resource, tspReqParams, Timeout)

	if tspErr != nil {
		return nil, newTransportErr(tspErr)
	}

	err = c.handleResponse(resBytes)

	return resBytes, err
}

func (c *Client) handleResponse(resBytes []byte) (err *Error) {
	// @TODO Refactor this conditions
	if c.tsp.IsWS() {
		res := &response{}
		parseErr := json.Unmarshal(resBytes, res)

		if parseErr != nil {
			return newJSONErr()
		}

		if res.Status == "error" {
			errMsg := strings.ToLower(res.Error)
			errCode := res.Code
			r := fmt.Sprintf("%d %s", errCode, errMsg)
			return &Error{name: ServiceErr, reason: r}
		}
	} else {
		if len(resBytes) == 0 {
			return nil
		}

		if isJSONArray(resBytes) {
			return nil
		}

		res := make(map[string]interface{})
		parseErr := json.Unmarshal(resBytes, &res)

		if parseErr != nil {
			return &Error{name: InvalidResponseErr, reason: parseErr.Error()}
		}

		if _, ok := res["message"]; !ok {
			return nil
		}

		httpRes := &httpResponse{
			Message: res["message"].(string),
		}
		if e, ok := res["error"].(float64); ok {
			httpRes.Status = int(e)
		} else {
			httpRes.Status = int(res["status"].(float64))
		}

		if httpRes.Status >= 400 {
			errMsg := strings.ToLower(httpRes.Message)
			errCode := httpRes.Status
			r := fmt.Sprintf("%d %s", errCode, errMsg)
			return &Error{name: ServiceErr, reason: r}
		}
	}

	return nil
}

func (c *Client) prepareRequestData(resourceName string, data map[string]interface{}) (resource string, reqParams *transport.RequestParams) {
	resource, method := c.resolveResource(resourceName, data)
	reqData := c.buildRequestData(resourceName, data)
	reqParams = c.createRequestParams(method, reqData)

	return resource, reqParams
}

func (c *Client) createRequestParams(method string, data interface{}) *transport.RequestParams {
	tspReqParams := &transport.RequestParams{
		Data: data,
	}

	if c.tsp.IsHTTP() {
		if method != "" {
			tspReqParams.Method = method
		}

		tspReqParams.AccessToken = c.accessToken
		tspReqParams.WaitTimeoutSeconds = c.PollingWaitTimeoutSeconds
	}

	return tspReqParams
}

func isJSONArray(b []byte) bool {
	return json.Unmarshal(b, &[]interface{}{}) == nil
}
