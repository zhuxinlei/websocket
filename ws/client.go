package ws

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Minute
	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type ClientKey string

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	key       ClientKey
	lastReqTs int64 // millisecond
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send  chan []byte
	close chan struct{}
}

// Failed send failed response to client.
func (c *Client) Failed(errorMsg string) {
	resp, _ := json.Marshal(SubscribeResponse{
		Status:   SubscribeStatusFailed,
		Ts:       c.lastReqTs,
		ErrorMsg: errorMsg,
	})
	c.send <- resp
}

// OK send ok response to client.
func (c *Client) OK(id string, sub, unsub Topic) {
	resp, _ := json.Marshal(SubscribeResponse{
		ID:       id,
		Status:   SubscribeStatusOk,
		Subbed:   sub,
		Unsubbed: unsub,
		Ts:       c.lastReqTs,
	})
	c.send <- resp
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		hubManager.unregister <- c
		c.conn.Close()
		close(c.close)
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(a string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		msgType, msg, err := c.conn.ReadMessage()
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// handle sub request.
		if msgType == websocket.TextMessage {
			// validate freq
			ts := time.Now().UnixNano() / 1000000
			if ts-c.lastReqTs < 100 {
				c.lastReqTs = ts
				c.Failed("request too frequently")
				continue
			}
			c.lastReqTs = ts

			// resolve req
			sub := SubscribeRequest{}
			err := json.Unmarshal(msg, &sub)
			if err != nil {
				c.Failed("bad request")
				continue
			}

			// sub case
			if sub.Sub != "" {
				if sub.Unsub != "" {
					c.Failed("both sub and unsub specified")
					continue
				}
				// validate topic
				if !sub.Sub.IsValid() {
					//验证topic是否动态类型
					if !sub.Sub.IsDynamic() {
						c.Failed("invalid topic")
						continue
					}
				}

				// sub
				done := make(chan struct{})
				hubManager.sub <- subscription{
					client: c,
					topic:  sub.Sub,
					done:   done,
				}
				<-done
				c.OK(sub.ID, sub.Sub, sub.Unsub)
				continue
			}

			// unsub case
			if sub.Unsub != "" {
				// validate topic
				if !sub.Unsub.IsValid() {
					c.Failed("invalid topic")
					continue
				}
				// unsub
				done := make(chan struct{})
				hubManager.unsub <- subscription{
					client: c,
					topic:  sub.Unsub,
					done:   done,
				}
				<-done
				c.OK(sub.ID, sub.Sub, sub.Unsub)
				continue
			}

		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			tryWriteMessage(w, message, 0)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Debugf("[%v] <- ping...", c.key)
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Debugf("[%v] <- ping failed", c.key)
				go func() {
					for i := 0; i < 3; i++ {
						time.Sleep(time.Second)
						c.conn.SetWriteDeadline(time.Now().Add(writeWait))
						if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
							log.Debugf("[%v] <- ping failed %d times", c.key, i+1)
							continue
						}
						return
					}
				}()
			}
		case <-c.close:
			return
		}
	}
}

func tryWriteMessage(writer io.WriteCloser, message []byte, count int) {
	if count == 5 {
		return
	}
	count++
	_, err := writer.Write(message)
	if err != nil {
		time.Sleep(time.Millisecond * 200)
		tryWriteMessage(writer, message, count)
	}
	return
}
