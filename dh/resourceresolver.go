package dh

func (c *Client) resolveResource(resourceName string) (resource, method string) {
	if c.tsp.IsHTTP() {
		rsrc, ok := httpResources[resourceName]

		if !ok {
			return resourceName, ""
		}

		return rsrc[0], rsrc[1]
	}

	if wsResources[resourceName] == "" {
		return resourceName, ""
	}

	return wsResources[resourceName], ""
}

var wsResources = map[string]string{
	"auth":        "authenticate",
	"tokenCreate": "token/create",
	"tokenRefresh": "token/refresh",
	"tokenByCreds": "token",
	"apiInfo": "server/info",
	"apiInfoCluster": "cluster/info",
}

var httpResources = map[string][2]string{
	"tokenCreate": [2]string{ "token/create", "POST" },
	"tokenRefresh": [2]string{ "token/refresh", "POST" },
	"tokenByCreds": [2]string{ "token", "POST" },
	"apiInfo": [2]string{ "info" },
	"apiInfoCluster": [2]string{ "info/config/cluster" },
}
