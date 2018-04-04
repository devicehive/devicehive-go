package dh

import (
	"encoding/json"
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
)

type Client struct {
	tsp transport.Transporter
}

func (c *Client) Authenticate(token string) (result bool, err *Error) {
	res, _, err := c.request(map[string]interface{}{
		"action": "authenticate",
		"token":  token,
	})

	if err != nil {
		return false, err
	}

	return res.Status == "success", nil
}

func (c *Client) request(data map[string]interface{}) (res *response, resBytes []byte, err *Error) {
	resBytes, tspErr := c.tsp.Request(data, Timeout)
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
