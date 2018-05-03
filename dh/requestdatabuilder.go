package dh

func (c *Client) buildRequestData(resourceName string, rawData map[string]interface{}) map[string]interface{} {
	if c.tsp.IsWS() {
		return rawData
	}

	payloadBuilder, ok := httpRequestPayloadBuilder[resourceName]

	if ok {
		return payloadBuilder(rawData)
	}

	return rawData
}

var httpRequestPayloadBuilder = map[string]func(map[string]interface{}) map[string]interface{} {
	"tokenCreate": func(data map[string]interface{}) map[string]interface{} {
		payload, ok := data["payload"].(map[string]interface{})

		if ok {
			return payload
		}

		return nil
	},
	"putConfig": func(data map[string]interface{}) map[string]interface{} {
		return map[string]interface{} {
			"value": data["value"],
		}
	},
	"deleteConfig": func(data map[string]interface{}) map[string]interface{} {
		return nil
	},
}
