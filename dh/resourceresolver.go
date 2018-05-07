package dh

import (
	"bytes"
	"text/template"
)

func (c *Client) resolveResource(resourceName string, data map[string]interface{}) (resource, method string) {
	if c.tsp.IsHTTP() {
		rsrc, ok := httpResources[resourceName]

		if !ok {
			return resourceName, ""
		}

		resource := prepareHttpResource(rsrc[0], data)
		method := rsrc[1]

		return resource, method
	}

	if wsResources[resourceName] == "" {
		return resourceName, ""
	}

	return wsResources[resourceName], ""
}

func prepareHttpResource(resourceTemplate string, data map[string]interface{}) string {
	t, err := template.New("resource").Parse(resourceTemplate)
	if err != nil {
		return ""
	}

	var resource bytes.Buffer
	err = t.Execute(&resource, data)
	if err != nil {
		return ""
	}

	return resource.String()
}

var wsResources = map[string]string{
	"auth":               "authenticate",
	"tokenCreate":        "token/create",
	"tokenRefresh":       "token/refresh",
	"tokenByCreds":       "token",
	"apiInfo":            "server/info",
	"apiInfoCluster":     "cluster/info",
	"putConfig":          "configuration/put",
	"getConfig":          "configuration/get",
	"deleteConfig":       "configuration/delete",
	"putDevice":          "device/save",
	"getDevice":          "device/get",
	"deleteDevice":       "device/delete",
	"insertCommand":      "command/insert",
	"listCommands":       "command/list",
	"updateCommand":      "command/update",
	"insertNotification": "notification/insert",
	"listNotifications":  "notification/list",
}

var httpResources = map[string][2]string{
	"tokenCreate":    {"token/create", "POST"},
	"tokenRefresh":   {"token/refresh", "POST"},
	"tokenByCreds":   {"token", "POST"},
	"apiInfo":        {"info"},
	"apiInfoCluster": {"info/config/cluster"},
	"putConfig":      {"configuration/{{.name}}", "PUT"},
	"getConfig":      {"configuration/{{.name}}"},
	"deleteConfig":   {"configuration/{{.name}}", "DELETE"},
	"putDevice":      {"device/{{.deviceId}}", "PUT"},
	"getDevice":      {"device/{{.deviceId}}"},
	"deleteDevice":   {"device/{{.deviceId}}", "DELETE"},
	"insertCommand":  {"device/{{.deviceId}}/command", "POST"},
	"listCommands": {
		`device/{{.deviceId}}/command?start={{or .start ""}}&end={{or .end ""}}&command={{or .command ""}}&status={{or .status ""}}&sortField={{or .sortField ""}}&sortOrder={{or .sortOrder ""}}&take={{or .take ""}}&skip={{or .skip ""}}`,
	},
	"updateCommand":      {"device/{{.deviceId}}/command/{{.commandId}}", "PUT"},
	"insertNotification": {"device/{{.deviceId}}/notification", "POST"},
	"listNotifications": {
		`device/{{.deviceId}}/notification?start={{or .start ""}}&end={{or .end ""}}&notification={{or .notification ""}}&sortField={{or .sortField ""}}&sortOrder={{or .sortOrder ""}}&take={{or .take ""}}&skip={{or .skip ""}}`,
	},
}
