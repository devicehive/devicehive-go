package dh

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"text/template"
)

func (c *Client) resolveResource(resourceName string, data map[string]interface{}) (resource, method string) {
	if c.tsp.IsHTTP() {
		rsrc, ok := httpResources[resourceName]

		if !ok {
			return resourceName, ""
		}

		queryParams := prepareQueryParams(data)

		resource := prepareHttpResource(rsrc[0], queryParams)
		method := rsrc[1]

		queryString := createQueryString(resourceName, queryParams)
		if queryString != "" {
			resource += "?" + queryString
		}

		return resource, method
	}

	if wsResources[resourceName] == "" {
		return resourceName, ""
	}

	return wsResources[resourceName], ""
}

func prepareHttpResource(resourceTemplate string, queryParams map[string]interface{}) string {
	t := template.New("resource")

	t, err := t.Parse(resourceTemplate)
	if err != nil {
		log.Printf("Error while parsing template: %s", err)
		return ""
	}

	var resource bytes.Buffer
	err = t.Execute(&resource, queryParams)
	if err != nil {
		log.Printf("Error while executing template: %s", err)
		return ""
	}

	return resource.String()
}

func prepareQueryParams(data map[string]interface{}) map[string]interface{} {
	preparedData := make(map[string]interface{})

	for k, v := range data {
		if s, ok := v.(string); ok {
			preparedData[k] = url.QueryEscape(s)
		} else if s, ok := v.([]string); ok {
			preparedData[k] = url.QueryEscape(strings.Join(s, ","))
		} else {
			preparedData[k] = v
		}
	}

	return preparedData
}

func createQueryString(resourceName string, queryParams map[string]interface{}) string {
	var params []string
	paramNames, ok := httpResourcesQueryParams[resourceName]

	if !ok {
		return ""
	}

	for _, p := range paramNames {
		if paramVal, ok := queryParams[p]; ok {
			params = append(params, fmt.Sprintf("%s=%v", p, paramVal))
		}
	}

	return strings.Join(params, "&")
}

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
	"deleteDevice":           "device/delete",
	"insertCommand":          "command/insert",
	"listCommands":           "command/list",
	"updateCommand":          "command/update",
	"insertNotification":     "notification/insert",
	"listNotifications":      "notification/list",
	"subscribeCommands":      "command/subscribe",
	"subscribeNotifications": "notification/subscribe",
}

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
