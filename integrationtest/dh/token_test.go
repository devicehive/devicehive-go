// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_test

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	is := is.New(t)

	expiration := time.Now().Add(1 * time.Second)
	accTok, refTok, err := client.CreateToken(*userId, expiration, expiration, nil, nil, nil)
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(accTok != "")
	is.True(refTok != "")
}

func TestRefreshToken(t *testing.T) {
	is := is.New(t)

	accessToken, err := client.RefreshToken()
	if err != nil {
		t.Fatalf("%s: %v", err.Name(), err)
	}

	is.True(accessToken != "")
}
