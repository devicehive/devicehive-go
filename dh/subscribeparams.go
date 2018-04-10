package dh

import (
	"time"
	"github.com/devicehive/devicehive-go/internal/utils"
	"errors"
)

type SubscribeParams struct {
	Action		  string 	`json:"action"`
	DeviceId      string    `json:"deviceId,omitempty"`
	NetworkIds    []string  `json:"networkIds,omitempty"`
	DeviceTypeIds []string  `json:"deviceTypeIds,omitempty"`
	Names         []string  `json:"names,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
}

func (p *SubscribeParams) Map() (m map[string]interface{}, err error) {
	m = utils.StructToJSONMap(p)

	if m == nil {
		return nil, errors.New("invalid JSON representation of struct")
	}

	if p.Timestamp.Unix() < 0 {
		delete(m, "timestamp")
	} else {
		m["timestamp"] = p.Timestamp.Format(timestampLayout)
	}

	return m, nil
}
