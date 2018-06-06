// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

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
