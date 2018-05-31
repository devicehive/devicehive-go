package devicehive_go

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
