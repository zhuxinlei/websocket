package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	host := "localhost:8002"
	port := "/ws/await"
	client := "1"
	connectInfo := connectInfo("e88016de-e268-4dbb-8a14-2e78aec3bf97", "123456")
	Connect(host, port, client, connectInfo, 3000*time.Second)
	select {}
}

func connectInfo(userID, connID string) http.Header {
	header := http.Header{}
	header.Add("User-Id", userID)
	header.Add("Connection-Id", connID)
	return header
}

func Connect(host, path, client string, header http.Header, pingInterval time.Duration) {
	u := url.URL{Scheme: "ws", Host: host, Path: path, RawQuery: "client=" + client}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Printf("dail: %v", err)
		if resp != nil {
			bs, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(bs))
		}
		return
	}
	defer c.Close()
	c.SetPingHandler(func(appData string) error {
		return nil
	})

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println("client <-", string(msg))
		}
	}()

	c.WriteMessage(websocket.TextMessage, []byte(`{"id":"1","sub":"otc"}`))

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			c.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
