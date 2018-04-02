package dh

import (
	"time"
)

type ServerInfo struct {
	APIVersion      string
	ServerTimestamp time.Time
	RestServerURL   string
}

func (c *Client) ServerInfo() (info *ServerInfo, err *Error) {
	res, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "server/info",
	})

	if err = c.handleResponseError(res, tspErr); err != nil {
		return nil, err
	}

	rawInfo := res["info"].(map[string]interface{})

	ts, tserr := time.Parse(timestampLayout, rawInfo["serverTimestamp"].(string))

	if tserr != nil {
		return nil, &Error{name: InvalidResponseErr, reason: tserr.Error()}
	}

	return &ServerInfo{
		APIVersion:      rawInfo["apiVersion"].(string),
		ServerTimestamp: ts,
		RestServerURL:   rawInfo["restServerUrl"].(string),
	}, nil
}

func (c *Client) ClusterInfo() (bootstrapServers string, zookeperConnect string, err *Error) {
	res, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "cluster/info",
	})

	if err = c.handleResponseError(res, tspErr); err != nil {
		return "", "", err
	}

	rawInfo := res["clusterInfo"].(map[string]interface{})

	return rawInfo["bootstrap.servers"].(string), rawInfo["zookeeper.connect"].(string), nil
}
