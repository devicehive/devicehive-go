package testutils

import (
	dh "github.com/devicehive/devicehive-go"
	"testing"
)

func LogDHErr(t *testing.T, err *dh.Error) {
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}
}
