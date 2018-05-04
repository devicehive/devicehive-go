package dh

import (
	"text/template"
	"bytes"
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
	"auth":        "authenticate",
	"tokenCreate": "token/create",
	"tokenRefresh": "token/refresh",
	"tokenByCreds": "token",
	"apiInfo": "server/info",
	"apiInfoCluster": "cluster/info",
	"putConfig": "configuration/put",
	"getConfig": "configuration/get",
	"deleteConfig": "configuration/delete",
	"putDevice": "device/save",
	"getDevice": "device/get",
	"deleteDevice": "device/delete",
}

var httpResources = map[string][2]string{
	"tokenCreate": [2]string{ "token/create", "POST" },
	"tokenRefresh": [2]string{ "token/refresh", "POST" },
	"tokenByCreds": [2]string{ "token", "POST" },
	"apiInfo": [2]string{ "info" },
	"apiInfoCluster": [2]string{ "info/config/cluster" },
	"putConfig": [2]string{ "configuration/{{index . `name`}}", "PUT" },
	"getConfig": [2]string{ "configuration/{{index . `name`}}" },
	"deleteConfig": [2]string{ "configuration/{{index . `name`}}", "DELETE" },
	"putDevice": [2]string{ "device/{{index . `deviceId`}}", "PUT" },
	"getDevice": [2]string{ "device/{{index . `deviceId`}}" },
	"deleteDevice": [2]string{ "device/{{index . `deviceId`}}", "DELETE" },
}
