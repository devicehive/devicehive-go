// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package devicehive_go

type token struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
