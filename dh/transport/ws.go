package transport

import (
	"github.com/gorilla/websocket"
	"strconv"
	"time"
	"math/rand"
)

func newWS(conn *websocket.Conn) *ws {
	tsp := &ws{
		conn: conn,
		requests: make(map[string]chan response),
	}

	go tsp.handleResponses()

	return tsp
}

type ws struct {
	conn *websocket.Conn
	requests map[string]chan response
}

func (t *ws) Request(data request) (res response, err error) {
	reqId := t.requestId()
	data["requestId"] = reqId
	err = t.conn.WriteJSON(data)

	if err != nil {
		return nil, err
	}

	t.requests[reqId] = make(chan response)

	res = response{}
	for {
		err = t.conn.ReadJSON(&res)

		if err != nil {
			return nil, err
		}

		if res["requestId"] == reqId {
			break
		}
	}

	return res, nil
}

func (t *ws) requestId() string {
	return strconv.FormatUint(rand.Uint64(), 10) + strconv.FormatInt(time.Now().Unix(), 10)
}

func (t *ws) handleResponses() {

}