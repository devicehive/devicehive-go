package transport

import (
	"math/rand"
	"strconv"
	"time"
)

type devicehiveData map[string]interface{}

var ranGen = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func (d devicehiveData) requestId() string {
	if d == nil {
		d = make(map[string]interface{})
	}

	reqId, ok := d["requestId"].(string)

	if !ok {
		r := strconv.FormatUint(ranGen.Uint64(), 10)
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		reqId = r + ts

		d["requestId"] = reqId
	}

	return reqId
}
