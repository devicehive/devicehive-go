// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

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
