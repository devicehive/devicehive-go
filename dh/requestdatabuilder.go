package dh

func (c *Client) buildRequestData(resourceName string, rawData map[string]interface{}) map[string]interface{} {
	if c.tsp.IsWS() {
		return rawData
	}

	payloadKey := httpRequestPayload[resourceName]

	if data, ok := rawData[payloadKey].(map[string]interface{}); payloadKey == "" || !ok {
		return rawData
	} else {
		return data
	}
}

var httpRequestPayload = map[string]string {
	"tokenCreate": "payload",
}
