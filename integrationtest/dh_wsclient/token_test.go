// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_wsclient_test

import (
	"testing"
	"time"
)

func TestWSClientCreateToken(t *testing.T) {
	expiration := time.Now().Add(1 * time.Second)
	err := wsclient.CreateToken(*userId, expiration, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}

func TestWSClientAccessTokenByRefresh(t *testing.T) {
	err := wsclient.AccessTokenByRefresh(*refreshToken)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}
