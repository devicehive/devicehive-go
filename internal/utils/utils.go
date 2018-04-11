package utils

import (
	"fmt"
	"github.com/fatih/structs"
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
	descr := structs.New(s)
	descr.TagName = "json"

	return descr.Map()
}
