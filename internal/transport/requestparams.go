// Copyright 2018 DataArt. All rights reserved.
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/devicehive/devicehive-go/internal/utils"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type RequestParams struct {
	Data               interface{}
	Method             string
	RequestId          string
	AccessToken        string
	WaitTimeoutSeconds int
}

var ranGen = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
var randLocker sync.Mutex

func (p *RequestParams) requestId() string {
	reqId := p.RequestId

	if reqId == "" {
		randLocker.Lock()
		r := strconv.FormatUint(ranGen.Uint64(), 10)
		randLocker.Unlock()
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		reqId = r + ts

		p.RequestId = reqId
	}

	return reqId
}

func (p *RequestParams) mapData() map[string]interface{} {
	data := utils.StructToJSONMap(p.Data)

	if data == nil {
		return make(map[string]interface{})
	}

	return data
}
