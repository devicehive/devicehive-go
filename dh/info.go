package dh

import (
	"encoding/json"
)

type serverInfo struct {
	Value *ServerInfo `json:"info"`
}

type ServerInfo struct {
	APIVersion      string      `json:"apiVersion"`
	ServerTimestamp ISO8601Time `json:"serverTimestamp"`
	RestServerURL   string      `json:"restServerUrl"`
}

type clusterInfo struct {
	Value *ClusterInfo `json:"clusterInfo"`
}

type ClusterInfo struct {
	BootstrapServers string `json:"bootstrap.servers"`
	ZookeeperConnect string `json:"zookeeper.connect"`
}

func (c *Client) ServerInfo() (info *ServerInfo, err *Error) {
	_, resBytes, err := c.request(map[string]interface{}{
		"action": "server/info",
	})

	if err != nil {
		return nil, err
	}

	info = &ServerInfo{}
	srvInfo := &serverInfo{Value: info}
	parseErr := json.Unmarshal(resBytes, srvInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return srvInfo.Value, nil
}

func (c *Client) ClusterInfo() (info *ClusterInfo, err *Error) {
	_, resBytes, err := c.request(map[string]interface{}{
		"action": "cluster/info",
	})

	if err != nil {
		return nil, err
	}

	info = &ClusterInfo{}
	clustInfo := &clusterInfo{Value: info}
	parseErr := json.Unmarshal(resBytes, clustInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return clustInfo.Value, nil
}
