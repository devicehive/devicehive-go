package testutils

import (
	"github.com/devicehive/devicehive-go/dh"
	"testing"
)

func LogDHErr(t *testing.T, err *dh.Error) {
	if err != nil {
		t.Errorf("%s: %v", err.Name(), err)
	}
}
