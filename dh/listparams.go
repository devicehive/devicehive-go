package dh

import (
	"time"
	"reflect"
)

type ListParams struct {
	DeviceId string `json:"deviceId"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	Notification string `json:"notification"`
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
	Take int `json:"take"`
	Skip int `json:"skip"`
}

func (lr *ListParams) Map() map[string]interface{} {
	res := make(map[string]interface{})
	t := reflect.TypeOf(*lr)
	v := reflect.ValueOf(*lr)

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		vf := v.Field(i)

		key := tf.Tag.Get("json")
		val := vf.Interface()

		if key != "" {
			switch val.(type) {
			case time.Time:

				res[key] = val.(time.Time).Format(timestampLayout)
			default:
				res[key] = val
			}
		}
	}

	return res
}
