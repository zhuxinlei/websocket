package ws

import (
	"github.com/gorilla/websocket"
)

func Register(conn *websocket.Conn, uid, connID string) {
	client := &Client{
		key:   ClientKey(uid),
		conn:  conn,
		send:  make(chan []byte, 16),
		close: make(chan struct{}),
	}
	hubManager.register <- client

	go client.writePump()
	go client.readPump()
}
