package dh

func (c *Client) buildRequestData(resourceName string, rawData map[string]interface{}) interface{} {
	var builders map[string]func(map[string]interface{}) interface{}
	if c.tsp.IsWS() {
		builders = wsRequestPayloadBuilder
	} else {
		builders = httpRequestPayloadBuilders
	}

	payloadBuilder, ok := builders[resourceName]

	if ok {
		return payloadBuilder(rawData)
	}

	return rawData
}

var wsRequestPayloadBuilder = map[string]func(map[string]interface{}) interface{} {}

var httpRequestPayloadBuilders = map[string]func(map[string]interface{}) interface{} {
	"tokenCreate": func(data map[string]interface{}) interface{} {
		return data["payload"]
	},
	"putConfig": func(data map[string]interface{}) interface{} {
		return map[string]interface{} {
			"value": data["value"],
		}
	},
	"deleteConfig": func(data map[string]interface{}) interface{} {
		return nil
	},
	"getConfig": func(data map[string]interface{}) interface{} {
		return nil
	},
	"putDevice": func(data map[string]interface{}) interface{} {
		return data["device"]
	},
	"getDevice": func(data map[string]interface{}) interface{} {
		return nil
	},
	"insertCommand": func(data map[string]interface{}) interface{} {
		return data["command"]
	},
	"listCommands": func(data map[string]interface{}) interface{} {
		return nil
	},
	"updateCommand": func(data map[string]interface{}) interface{} {
		return data["command"]
	},
	"insertNotification": func(data map[string]interface{}) interface{} {
		return data["notification"]
	},
}
