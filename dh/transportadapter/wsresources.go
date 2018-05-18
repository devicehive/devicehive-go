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
}
