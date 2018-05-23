package dh

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
	info = &ServerInfo{}

	err = c.getModel("apiInfo", info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *Client) GetClusterInfo() (info *ClusterInfo, err *Error) {
	info = &ClusterInfo{}

	err = c.getModel("apiInfoCluster", info, nil)
	if err != nil {
		return nil, err
	}

	return info, nil
}
