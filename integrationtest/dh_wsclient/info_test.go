// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_wsclient_test

import (
	"testing"
)

func TestWSClientGetInfo(t *testing.T) {
	err := wsclient.GetInfo()
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestWSClientGetClusterInfo(t *testing.T) {
	err := wsclient.GetClusterInfo()
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}
