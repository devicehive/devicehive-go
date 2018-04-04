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
	resBytes, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "authenticate",
		"token":  token,
	})

	var res *response
	if res, err = c.handleResponse(resBytes, tspErr); err != nil {
		return false, err
	}

	return res.Status == "success", nil
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
