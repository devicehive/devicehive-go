// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"strconv"
	"strings"

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

func JoinIntSlice(i []int, sep string) string {
	return strings.Join(IntSliceToStrSlice(i), sep)
}

func IntSliceToStrSlice(s []int) (stringSlice []string) {
	for _, elem := range s {
		stringSlice = append(stringSlice, strconv.FormatInt(int64(elem), 10))
	}

	return stringSlice
}

func StructToJSONMap(s interface{}) map[string]interface{} {
	if s == nil {
		return nil
	}

	if m, ok := s.(map[string]interface{}); ok {
		return m
	}

	descr := structs.New(s)
	descr.TagName = "json"

	return descr.Map()
}
