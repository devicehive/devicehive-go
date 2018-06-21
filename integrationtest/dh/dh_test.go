// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_test

import (
	"flag"
	"fmt"
	dh "github.com/devicehive/devicehive-go"
	"os"
	"testing"
	"time"
)

var serverAddr = flag.String("serverAddress", "", "Server address without trailing slash")
var accessToken = flag.String("accessToken", "", "Your access token")
var refreshToken = flag.String("refreshToken", "", "Your refresh token")
var userId = flag.Int("userId", 0, "DH user ID")

var client *dh.Client

var waitTimeout time.Duration

func TestMain(m *testing.M) {
	flag.Parse()

	if *accessToken == "" || *refreshToken == "" || *userId == 0 {
		os.Exit(1)
	}

	var err *dh.Error
	client, err = dh.ConnectWithToken(*serverAddr, *accessToken, *refreshToken, nil)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client.PollingWaitTimeoutSeconds = 7

	waitTimeout = time.Duration(client.PollingWaitTimeoutSeconds+1) * time.Second

	res := m.Run()
	os.Exit(res)
}
