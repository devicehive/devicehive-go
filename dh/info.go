package dh

import (
	"encoding/json"
)

type ServerInfo struct {
	APIVersion         string      `json:"apiVersion"`
	ServerTimestamp    ISO8601Time `json:"serverTimestamp"`
	RestServerURL      string      `json:"restServerUrl"`
	WebSocketServerURL string      `json:"webSocketServerUrl"`
}

type ClusterInfo struct {
	BootstrapServers string `json:"bootstrap.servers"`
	ZookeeperConnect string `json:"zookeeper.connect"`
}

func (c *Client) GetInfo() (info *ServerInfo, err *Error) {
	rawRes, err := c.request("apiInfo", nil)

	if err != nil {
		return nil, err
	}

	info = &ServerInfo{}
	parseErr := json.Unmarshal(rawRes, info)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	return info, nil
}

func (c *Client) GetClusterInfo() (info *ClusterInfo, err *Error) {
	rawRes, err := c.request("apiInfoCluster", nil)

	if err != nil {
		return nil, err
	}

	info = &ClusterInfo{}
	parseErr := json.Unmarshal(rawRes, info)
	if parseErr != nil {
		return nil, newJSONErr()
	}

	return info, nil
}
