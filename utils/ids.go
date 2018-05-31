package utils

import "encoding/json"

func ParseIDs(b []byte) (ids *IDs, err error) {
	ids = &IDs{}
	err = json.Unmarshal(b, ids)
	return ids, err
}

type IDs struct {
	Request      string `json:"requestId"`
	Subscription int64  `json:"subscriptionId"`
}
