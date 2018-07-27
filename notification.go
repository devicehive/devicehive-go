// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

type Notification struct {
	Id           int                    `json:"id"`
	Notification string                 `json:"notification"`
	Timestamp    ISO8601Time            `json:"timestamp"`
	DeviceId     string                 `json:"deviceId"`
	NetworkId    int                    `json:"networkId"`
	Parameters   map[string]interface{} `json:"parameters"`
}
