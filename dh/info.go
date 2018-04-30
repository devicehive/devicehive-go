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

func (c *Client) GetInfo() (info *ServerInfo, err *Error) {
	rawRes, err := c.request("server/info", nil)

	if err != nil {
		return nil, err
	}

	info = &ServerInfo{}
	srvInfo := &serverInfo{Value: info}
	parseErr := json.Unmarshal(rawRes, srvInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return srvInfo.Value, nil
}

func (c *Client) GetClusterInfo() (info *ClusterInfo, err *Error) {
	rawRes, err := c.request("cluster/info", nil)

	if err != nil {
		return nil, err
	}

	info = &ClusterInfo{}
	clustInfo := &clusterInfo{Value: info}
	parseErr := json.Unmarshal(rawRes, clustInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return clustInfo.Value, nil
}
