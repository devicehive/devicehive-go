package dh

import (
	"strings"
	"time"
)

type dhTime struct {
	time.Time
}

func (t *dhTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	if s == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(timestampLayout, s)
	return err
}
