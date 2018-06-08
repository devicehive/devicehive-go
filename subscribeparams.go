// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"errors"
	"time"

	"github.com/devicehive/devicehive-go/utils"
)

type SubscribeParams struct {
	DeviceId              string    `json:"deviceId,omitempty"`
	NetworkIds            []int     `json:"networkIds,omitempty"`
	DeviceTypeIds         []int     `json:"deviceTypeIds,omitempty"`
	Names                 []string  `json:"names,omitempty"`
	Timestamp             time.Time `json:"timestamp,omitempty"`
	ReturnUpdatedCommands bool      `json:"returnUpdatedCommands,omitempty"`
	Limit                 int       `json:"limit,omitempty"`
	WaitTimeout           int       `json:"waitTimeout,omitempty"`
}

func (p *SubscribeParams) Map() (map[string]interface{}, error) {
	m := utils.StructToJSONMap(p)

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
