// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package dh_wsclient_test

import (
	"encoding/json"
	"github.com/matryer/is"
	"testing"
)

func TestProperty(t *testing.T) {
	const testPropertyName = "go-test-prop"
	is := is.New(t)

	err := wsclient.SetProperty(testPropertyName, "go-test-val")
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		is.Equal(res["name"].(string), testPropertyName)
	})

	err = wsclient.GetProperty(testPropertyName)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, func(data []byte) {
		res := make(map[string]interface{})
		json.Unmarshal(data, &res)
		is.Equal(res["value"].(string), "go-test-val")
	})

	err = wsclient.DeleteProperty(testPropertyName)
	if err != nil {
		t.Fatal(err)
	}

	testResponse(t, nil)
}
