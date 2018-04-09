package dh

import (
	"time"
	"encoding/json"
)

type ListParams struct {
	Action string `json:"action"`
	DeviceId string `json:"deviceId,omitempty"`
	Start time.Time `json:"start,omitempty"`
	End time.Time `json:"end,omitempty"`
	Notification string `json:"notification,omitempty"`
	SortField string `json:"sortField,omitempty"`
	SortOrder string `json:"sortOrder,omitempty"`
	Take int `json:"take,omitempty"`
	Skip int `json:"skip,omitempty"`
}

func (p *ListParams) Map() (res map[string]interface{}, err error) {
	res = make(map[string]interface{})

	b, err := json.Marshal(p)

	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(b, &res)

	if p.Start.Unix() < 0 {
		delete(res, "start")
	} else {
		res["start"] = p.Start.Format(timestampLayout)
	}
	if p.End.Unix() < 0 {
		delete(res, "end")
	} else {
		res["end"] = p.End.Format(timestampLayout)
	}

	return res, nil
}
