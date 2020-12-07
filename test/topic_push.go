package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var WSCli = &WSClient{
	Client: http.Client{
		Timeout: 15 * time.Second,
	},
}

type WSClient struct {
	http.Client
}

type WsPushData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type TopicPushReq struct {
	Targets []string    `json:"targets"`
	Content interface{} `json:"content"`
}

func TopicPush(topic string, targets []string, message WsPushData) {
	request := TopicPushReq{}
	request.Targets = targets
	request.Content = message
	bs, err := json.Marshal(request)
	if err != nil {
		log.Errorf("ws.TopicPush marshal %+v failed: %v", request, err)
		return
	}
	url := fmt.Sprintf("http://localhost:8002/ws/push/topic/%s", topic)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
	if err != nil {
		log.Errorf("topic push new req failed: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := WSCli.Do(req)
	if err != nil {
		log.Errorf("topic push do req failed: %v", err)
		return
	}
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}
