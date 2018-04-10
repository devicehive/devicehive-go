package utils

import (
	"fmt"
	"encoding/json"
)

func ISliceToStrSlice(s []interface{}) (stringSlice []string, err error) {
	for _, elem := range s {
		switch elem.(type) {
		case string:
			stringSlice = append(stringSlice, elem.(string))
		default:
			return nil, fmt.Errorf("element is not string: %v", elem)
		}
	}

	return stringSlice, nil
}

func StructToJSONMap(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})

	b, err := json.Marshal(s)

	if err != nil {
		return nil
	}

	_ = json.Unmarshal(b, &m)

	return m
}