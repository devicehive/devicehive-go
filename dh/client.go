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
		res, _, err := c.request("auth", map[string]interface{}{
			"token": token,
		})

		if err != nil {
			return false, err
		}

		return res.Status == "success", nil
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

	_, rawRes, err := c.request(resourceName, data)

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
	_, _, err := c.request(resourceName, map[string]interface{}{
		"subscriptionId": subscriptionId,
	})

	if err != nil {
		return err
	}

	c.tsp.Unsubscribe(subscriptionId)

	return nil
}

func (c *Client) request(resourceName string, data map[string]interface{}) (res *response, resBytes []byte, err *Error) {
	resource, method := c.resolveResource(resourceName)

	if resource == "" {
		return nil, nil, &Error{name: InvalidRequestErr, reason: "unknown resource name"}
	}

	tspReqParams := &transport.RequestParams{
		Data: make(map[string]interface{}),
	}

	if c.tsp.IsHTTP() && method != "" {
		tspReqParams.Method = method
	}

	for k, v := range data {
		tspReqParams.Data[k] = v
	}

	resBytes, tspErr := c.tsp.Request(resource, tspReqParams, Timeout)
	res, err = c.handleResponse(resBytes, tspErr)

	return res, resBytes, err
}

func (c *Client) handleResponse(resBytes []byte, tspErr *transport.Error) (res *response, err *Error) {
	if tspErr != nil {
		return nil, newTransportErr(tspErr)
	}

	res = &response{}
	parseErr := json.Unmarshal(resBytes, res)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	if res.Status == "error" {
		errMsg := strings.ToLower(res.Error)
		errCode := res.Code
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return nil, &Error{name: ServiceErr, reason: r}
	}

	return res, nil
}
