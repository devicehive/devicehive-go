package httputils

import (
	"log"
	"bytes"
	"github.com/devicehive/devicehive-go/internal/utils"
	"text/template"
	"fmt"
	"net/url"
	"strings"
)

func PrepareHttpResource(resourceTemplate string, queryParams map[string]string) string {
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

func PrepareQueryParams(data map[string]interface{}) map[string]string {
	preparedData := make(map[string]string)

	for k, v := range data {
		if s, ok := v.(string); ok {
			preparedData[k] = url.QueryEscape(s)
		} else if s, ok := v.([]string); ok {
			preparedData[k] = url.QueryEscape(strings.Join(s, ","))
		} else if i, ok := v.([]int); ok {
			preparedData[k] = url.QueryEscape(utils.JoinIntSlice(i, ","))
		} else {
			preparedData[k] = fmt.Sprintf("%v", v)
		}
	}

	return preparedData
}

func CreateQueryString(resourcesQueryParams map[string][]string, resourceName string, queryParams map[string]string) string {
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
