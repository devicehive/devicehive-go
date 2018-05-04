package dh

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strconv"
	"strings"
)

type Client struct {
	tsp          transport.Transporter
	accessToken  string
	refreshToken string
	login        string
	password     string
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

	rawRes, err := c.request(resourceName, data)

	if err != nil {
		return nil, "", err
	}

	type subsId struct {
		Value int64 `json:"subscriptionId"`
	}
	id := &subsId{}

	parseErr := json.Unmarshal(rawRes, id)

	if parseErr != nil {
		return nil, "", newJSONErr()
	}

	subscriptionId = strconv.FormatInt(id.Value, 10)

	return c.tsp.Subscribe(subscriptionId), subscriptionId, nil
}

func (c *Client) unsubscribe(resourceName, subscriptionId string) *Error {
	_, err := c.request(resourceName, map[string]interface{}{
		"subscriptionId": subscriptionId,
	})

	if err != nil {
		return err
	}

	c.tsp.Unsubscribe(subscriptionId)

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (resBytes []byte, err *Error) {
	resource, method := c.resolveResource(resourceName, data)

	if resource == "" {
		return nil, &Error{name: InvalidRequestErr, reason: "unknown resource name"}
	}

	reqData := c.buildRequestData(resourceName, data)
	tspReqParams := c.createRequestParams(method, reqData)

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

		res := &httpResponse{}

		parseErr := json.Unmarshal(resBytes, res)

		if parseErr != nil {
			return newJSONErr()
		}

		if res.Status >= 400 {
			errMsg := strings.ToLower(res.Error)
			errCode := res.Status
			r := fmt.Sprintf("%d %s", errCode, errMsg)
			return &Error{name: ServiceErr, reason: r}
		}
	}

	return nil
}

func (c *Client) createRequestParams(method string, data map[string]interface{}) *transport.RequestParams {
	tspReqParams := &transport.RequestParams{
		Data: make(map[string]interface{}),
	}

	if c.tsp.IsHTTP() {
		if method != "" {
			tspReqParams.Method = method
		}

		tspReqParams.AccessToken = c.accessToken
	}

	for k, v := range data {
		tspReqParams.Data[k] = v
	}

	return tspReqParams
}
