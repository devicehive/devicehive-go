package dh

import (
	"time"
	"encoding/json"
)

type serverInfo struct {
	Value *ServerInfo `json:"info"`
}

type ServerInfo struct {
	APIVersion      string `json:"apiVersion"`
	ServerTimestamp time.Time `json:"serverTimestamp"`
	RestServerURL   string `json:"restServerUrl"`
}

type clusterInfo struct {
	Value *ClusterInfo `json:"clusterInfo"`
}

type ClusterInfo struct {
	BootstrapServers string `json:"bootstrap.servers"`
	ZookeeperConnect string `json:"zookeeper.connect"`
}

func (c *Client) ServerInfo() (info *ServerInfo, err *Error) {
	resBytes, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "server/info",
	})

	if _, err = c.handleResponse(resBytes, tspErr); err != nil {
		return nil, err
	}

	info = &ServerInfo{}
	srvInfo := &serverInfo{ Value: info }
	parseErr := json.Unmarshal(resBytes, srvInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	//ts, tserr := time.Parse(timestampLayout, rawInfo["serverTimestamp"].(string))
	//
	//if tserr != nil {
	//	return nil, &Error{name: InvalidResponseErr, reason: tserr.Error()}
	//}

	return srvInfo.Value, nil
}

func (c *Client) ClusterInfo() (info *ClusterInfo, err *Error) {
	resBytes, tspErr := c.tsp.Request(map[string]interface{}{
		"action": "cluster/info",
	})

	if _, err = c.handleResponse(resBytes, tspErr); err != nil {
		return nil, err
	}

	info = &ClusterInfo{}
	clustInfo := &clusterInfo{ Value: info }
	parseErr := json.Unmarshal(resBytes, clustInfo)

	if parseErr != nil {
		return nil, newJSONErr()
	}

	return clustInfo.Value, nil
}
