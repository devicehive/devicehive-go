package dh

import (
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"strconv"
)

type dhClient struct {
	conn *websocket.Conn
}

func (c *dhClient) Authenticate(token string) (result bool, err error) {
	rand.Seed(time.Now().Unix())
	reqId := strconv.Itoa(rand.Int())

	err = c.conn.WriteJSON(map[string]string{
		"action": "authenticate",
		"requestId": reqId,
		"token": token,
	})

	if err != nil {
		return false, err
	}

	res := make(map[string]interface{})
	for {
		err = c.conn.ReadJSON(&res)

		if err != nil {
			return false, err
		}

		if res["requestId"] == reqId {
			break
		}
	}

	return res["status"] == "success", nil
}