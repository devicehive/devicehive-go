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
	"listDevices":            {"device"},
	"deleteDevice":           {"device/{{.deviceId}}", "DELETE"},
	"insertCommand":          {"device/{{.deviceId}}/command", "POST"},
	"listCommands":           {`device/{{.deviceId}}/command`},
	"updateCommand":          {"device/{{.deviceId}}/command/{{.commandId}}", "PUT"},
	"insertNotification":     {"device/{{.deviceId}}/notification", "POST"},
	"listNotifications":      {"device/{{.deviceId}}/notification"},
	"subscribeCommands":      {"device/command/poll"},
	"subscribeNotifications": {"device/notification/poll"},
	"insertNetwork":          {"network", "POST"},
	"deleteNetwork":          {"network/{{.networkId}}", "DELETE"},
	"updateNetwork":          {"network/{{.networkId}}", "PUT"},
	"getNetwork":             {"network/{{.networkId}}"},
	"listNetworks":           {"network"},
	"insertDeviceType":       {"devicetype", "POST"},
	"updateDeviceType":       {"devicetype/{{.deviceTypeId}}", "PUT"},
	"deleteDeviceType":       {"devicetype/{{.deviceTypeId}}", "DELETE"},
	"getDeviceType":          {"devicetype/{{.deviceTypeId}}"},
	"listDeviceTypes":        {"devicetype"},
	"createUser":             {"user", "POST"},
	"deleteUser":             {"user/{{.userId}}", "DELETE"},
	"getUser":                {"user/{{.userId}}"},
	"getCurrentUser":         {"user/current"},
	"listUsers":              {"user"},
	"updateUser":			  {"user/{{.userId}}", "PUT"},
	"assignNetwork":		  {"user/{{.userId}}/network/{{.networkId}}", "PUT"},
}

var httpResourcesQueryParams = map[string][]string{
	"listCommands":           {"start", "end", "command", "status", "sortField", "sortOrder", "take", "skip"},
	"listNotifications":      {"start", "end", "notification", "sortField", "sortOrder", "take", "skip"},
	"subscribeCommands":      {"deviceId", "timestamp", "waitTimeout", "names"},
	"subscribeNotifications": {"deviceId", "timestamp", "waitTimeout", "names"},
	"listDevices":            {"name", "namePattern", "networkId", "networkName", "sortField", "sortOrder", "take", "skip"},
	"listNetworks":           {"name", "namePattern", "sortField", "sortOrder", "take", "skip"},
	"listDeviceTypes":        {"name", "namePattern", "sortField", "sortOrder", "take", "skip"},
	"listUsers":              {"login", "loginPattern", "role", "status", "sortField", "sortOrder", "take", "skip"},
}
