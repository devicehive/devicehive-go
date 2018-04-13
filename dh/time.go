package dh

import (
	"strings"
	"time"
)

const (
	timestampLayout  = "2006-01-02T15:04:05.000"
)

type ISO8601Time struct {
	time.Time
}

func (t *ISO8601Time) FormatISO8601() string {
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
