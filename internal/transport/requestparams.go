package transport

import (
	"math/rand"
	"strconv"
	"time"
)

type RequestParams struct {
	Action    string
	Data      map[string]interface{}
	Method    string
	RequestId string
}

var ranGen = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func (p *RequestParams) requestId() string {
	reqId := p.RequestId

	if reqId == "" {
		r := strconv.FormatUint(ranGen.Uint64(), 10)
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		reqId = r + ts

		p.RequestId = reqId
	}

	return reqId
}
