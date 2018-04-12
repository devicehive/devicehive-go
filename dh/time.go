package dh

import (
	"strings"
	"time"
)

type ISO8601Time struct {
	*time.Time
}

func (t *ISO8601Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	if s == "null" || s == "" {
		t.Time = &time.Time{}
		return
	}

	parsedTime, err := time.Parse(timestampLayout, s)

	if err != nil {
		return err
	}

	t.Time = &parsedTime

	return nil
}
