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
	refreshToken string
	login        string
	password     string
}

func (c *Client) Authenticate(token string) (result bool, err *Error) {
	res, _, err := c.request("authenticate", map[string]interface{}{
		"token":  token,
	})

	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}

func (c *Client) subscribe(resource string, params *SubscribeParams) (tspChan chan []byte, subscriptionId string, err *Error) {
	if params == nil {
		params = &SubscribeParams{}
	}

	data, jsonErr := params.Map()

	if jsonErr != nil {
		return nil, "", &Error{name: InvalidRequestErr, reason: jsonErr.Error()}
	}

	_, rawRes, err := c.request(resource, data)

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

func (c *Client) unsubscribe(resource, subscriptionId string) *Error {
	_, _, err := c.request(resource, map[string]interface{}{
		"subscriptionId": subscriptionId,
	})

	if err != nil {
		return err
	}

	c.tsp.Unsubscribe(subscriptionId)

	return nil
}

func (c *Client) request(resource string, data map[string]interface{}) (res *response, resBytes []byte, err *Error) {
	resBytes, tspErr := c.tsp.Request(resource, data, Timeout)
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
