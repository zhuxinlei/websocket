package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"server-tokenhouse-ws/service/respond"
	"server-tokenhouse-ws/ws"
)

func SinglePush(c *gin.Context) {
	userID := c.Param("uid")
	//connID := c.Param("cid")

	message, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "general", err.Error())
		return
	}

	ws.GetHubManager().SinglePush <- &ws.SingleMessage{
		Target:  ws.ClientKey(userID),
		Content: message,
	}

	respond.EmptySuccess(c)
}

func Broadcast(c *gin.Context) {

	message, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "general", err.Error())
		return
	}

	ws.GetHubManager().BroadcastPush <- message

	respond.EmptySuccess(c)
}

type TopicPushReq struct {
	Targets []ws.ClientKey  `json:"targets"`
	Content json.RawMessage `json:"content"`
}

func Topic(c *gin.Context) {

	topic := c.Param("topic")

	req := TopicPushReq{}
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "general", err.Error())
		return
	}

	ws.GetHubManager().TopicPush <- &ws.TopicMessage{
		Topic:   ws.Topic(topic),
		Targets: req.Targets,
		Content: req.Content,
	}

	respond.EmptySuccess(c)
}

func Topic2(c *gin.Context) {

	topic := c.Param("topic")
	uniqueKey := c.Param("unique_key")
	topicName := topic+"/"+uniqueKey
	req := TopicPushReq{}
	if err := c.ShouldBind(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "general", err.Error())
		return
	}

	ws.GetHubManager().TopicPush <- &ws.TopicMessage{
		Topic:   ws.Topic(topicName),
		Targets: req.Targets,
		Content: req.Content,
	}

	respond.EmptySuccess(c)
}
