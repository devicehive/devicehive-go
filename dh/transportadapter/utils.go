package transportadapter

import (
	"encoding/json"
	"log"
	"bytes"
	"net/url"
	"strings"
	"fmt"
	"text/template"
)

func isJSONArray(b []byte) bool {
	return json.Unmarshal(b, &[]interface{}{}) == nil
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

func createQueryString(resourcesQueryParams map[string][]string, resourceName string, queryParams map[string]interface{}) string {
	var params []string
	paramNames, ok := resourcesQueryParams[resourceName]

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
