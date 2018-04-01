package dh

import (
	"fmt"
	"github.com/devicehive/devicehive-go/internal/transport"
	"strings"
)

type Client struct {
	tsp transport.Transporter
}

func (c *Client) handleResponseError(response map[string]interface{}, err *transport.Error) *Error {
	if err != nil {
		return newTransportErr(err)
	}

	if response["status"] == "error" {
		errMsg := strings.ToLower(response["error"].(string))
		errCode := int(response["code"].(float64))
		r := fmt.Sprintf("%d %s", errCode, errMsg)
		return &Error{name: ServiceErr, reason: r}
	}

	return nil
}

func (c *Client) Authenticate(token string) (result bool, err *Error) {
	res, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "authenticate",
		"token":  token,
	})

	if err = c.handleResponseError(res, tspErr); err != nil {
		return false, err
	}

	return res["status"].(string) == "success", nil
}