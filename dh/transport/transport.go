package transport

import (
	"github.com/gorilla/websocket"
	"strconv"
	"time"
	"math/rand"
)

type Transporter interface {
	Request(data devicehiveData) (res devicehiveData, err error)
}

func Create(url string) (transport Transporter, err error) {
	// @TODO HTTP transport
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		return nil, err
	}

	return newWS(conn), nil
}

type devicehiveData map[string]interface{}

var ranGen = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
func (d devicehiveData) requestId() string {
	reqId, ok := d["requestId"].(string)

	if !ok {
		r := strconv.FormatUint(ranGen.Uint64(), 10)
		ts := strconv.FormatInt(time.Now().Unix(), 10)
		reqId = r + ts

		d["requestId"] = reqId
	}

	return reqId
}