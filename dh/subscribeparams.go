package dh

import (
	"time"
	"encoding/json"
)

type SubscribeParams struct {
	Action		  string 	`json:"action"`
	DeviceId      string    `json:"deviceId"`
	NetworkIds    []string  `json:"networkIds"`
	DeviceTypeIds []string  `json:"deviceTypeIds"`
	Names         []string  `json:"names"`
	Timestamp     time.Time `json:"timestamp"`
}

func (p *SubscribeParams) Map() (res map[string]interface{}, err error) {
	res = make(map[string]interface{})

	b, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(b, &res)

	if p.Timestamp.Unix() < 0 {
		delete(res, "timestamp")
	} else {
		res["timestamp"] = p.Timestamp.Format(timestampLayout)
	}

	return res, nil
}
