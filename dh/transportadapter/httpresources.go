package transportadapter

var httpResources = map[string][2]string{
	"tokenCreate":            {"token/create", "POST"},
	"tokenRefresh":           {"token/refresh", "POST"},
	"tokenByCreds":           {"token", "POST"},
	"apiInfo":                {"info"},
	"apiInfoCluster":         {"info/config/cluster"},
	"putConfig":              {"configuration/{{.name}}", "PUT"},
	"getConfig":              {"configuration/{{.name}}"},
	"deleteConfig":           {"configuration/{{.name}}", "DELETE"},
	"putDevice":              {"device/{{.deviceId}}", "PUT"},
	"getDevice":              {"device/{{.deviceId}}"},
	"deleteDevice":           {"device/{{.deviceId}}", "DELETE"},
	"insertCommand":          {"device/{{.deviceId}}/command", "POST"},
	"listCommands":           {`device/{{.deviceId}}/command`},
	"updateCommand":          {"device/{{.deviceId}}/command/{{.commandId}}", "PUT"},
	"insertNotification":     {"device/{{.deviceId}}/notification", "POST"},
	"listNotifications":      {`device/{{.deviceId}}/notification`},
	"subscribeCommands":      {`device/command/poll`},
	"subscribeNotifications": {`device/notification/poll`},
}

var httpResourcesQueryParams = map[string][]string{
	"listCommands":           {"start", "end", "command", "status", "sortField", "sortOrder", "take", "skip"},
	"listNotifications":      {"start", "end", "notification", "sortField", "sortOrder", "take", "skip"},
	"subscribeCommands":      {"deviceId", "timestamp", "waitTimeout", "names"},
	"subscribeNotifications": {"deviceId", "timestamp", "waitTimeout", "names"},
}