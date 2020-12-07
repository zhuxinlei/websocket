package test

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"net/url"
	"testing"
)

func TestNoUnsub(t *testing.T) {
	connAndDisConn(true)
	connAndDisConn(true)
}

func connAndDisConn(sub bool) {
	u := url.URL{Scheme: "ws", Host: "localhost:8002", Path: "/ws/await"}

	header := http.Header{}
	header.Add("User-Id", uuid.NewV4().String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return
	}
	defer c.Close()

	if sub {
		c.WriteMessage(websocket.TextMessage, []byte(`{"id":"1","sub":"otc"}`))
	}

}
