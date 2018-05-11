package dh

import (
	"bytes"
	"text/template"
	"log"
	"strings"
	"net/url"
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
	t := template.New("resource")

	t, err := t.Funcs(template.FuncMap{"Join": func(s []string, sep string) string {
		str := strings.Join(s, sep)
		return url.QueryEscape(str)
	}}).Parse(resourceTemplate)
	if err != nil {
		log.Printf("Error while parsing template: %s", err)
		return ""
	}

	var resource bytes.Buffer
	err = t.Execute(&resource, data)
	if err != nil {
		log.Printf("Error while executing template: %s", err)
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
	"subscribeCommands":  "command/subscribe",
	"subscribeNotifications": "notification/subscribe",
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
	"subscribeCommands": {`device/command/poll?deviceId={{or .deviceId ""}}&networkIds={{Join .networkIds ","}}&deviceTypeIds={{Join .deviceTypeIds ","}}&names={{Join .names ","}}&timestamp={{or .timestamp ""}}&waitTimeout={{or .waitTimeout 0}}`},
	"subscribeNotifications": {`device/notification/poll?deviceId={{or .deviceId ""}}&networkIds={{Join .networkIds ","}}&deviceTypeIds={{Join .deviceTypeIds ","}}&names={{Join .names ","}}&timestamp={{or .timestamp ""}}&waitTimeout={{or .waitTimeout 0}}`},
}
