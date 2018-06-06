// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

import (
	"strings"
	"time"
)

// Custom timestamp in ISO8601 format
type ISO8601Time struct {
	time.Time
}

func (t *ISO8601Time) String() string {
	return t.Time.Format(timestampLayout)
}

func (t *ISO8601Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(timestampLayout, s)

	if err != nil {
		return err
	}

	return nil
}

func (t *ISO8601Time) MarshalJSON() (b []byte, err error) {
	if t.Time.UnixNano() <= 0 {
		return []byte("\"\""), nil
	}

	return []byte("\"" + t.String() + "\""), nil
}
