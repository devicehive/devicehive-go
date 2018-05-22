package transportadapter

var wsResources = map[string]string{
	"auth":                   "authenticate",
	"tokenCreate":            "token/create",
	"tokenRefresh":           "token/refresh",
	"tokenByCreds":           "token",
	"apiInfo":                "server/info",
	"apiInfoCluster":         "cluster/info",
	"putConfig":              "configuration/put",
	"getConfig":              "configuration/get",
	"deleteConfig":           "configuration/delete",
	"putDevice":              "device/save",
	"getDevice":              "device/get",
	"listDevices":            "device/list",
	"deleteDevice":           "device/delete",
	"insertCommand":          "command/insert",
	"listCommands":           "command/list",
	"updateCommand":          "command/update",
	"insertNotification":     "notification/insert",
	"listNotifications":      "notification/list",
	"subscribeCommands":      "command/subscribe",
	"subscribeNotifications": "notification/subscribe",
	"insertNetwork":          "network/insert",
	"deleteNetwork":          "network/delete",
	"updateNetwork":          "network/update",
	"getNetwork":             "network/get",
	"listNetworks":           "network/list",
	"insertDeviceType":       "devicetype/insert",
	"updateDeviceType":       "devicetype/update",
	"deleteDeviceType":       "devicetype/delete",
	"getDeviceType":          "devicetype/get",
	"listDeviceTypes":        "devicetype/list",
	"createUser":             "user/insert",
	"deleteUser":             "user/delete",
	"getUser":                "user/get",
	"getCurrentUser":         "user/getCurrent",
	"listUsers":              "user/list",
}
