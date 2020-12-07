package service

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"server-tokenhouse-ws/service/handlers"
)

func Run(host string, port int) {
	router := gin.Default()

	// websocket connect
	router.GET("/ws/await", handlers.Await)
	// single push
	router.POST("/ws/push/single/:uid", handlers.SinglePush)
	// broadcast
	router.POST("/ws/push/broadcast", handlers.Broadcast)
	// topic
	router.POST("/ws/push/topic/:topic", handlers.Topic)

	router.POST("/ws/push/topic/:topic/:unique_key", handlers.Topic2)

	ginpprof.Wrap(router)

	router.Run(fmt.Sprintf("%s:%d", host, port))
}
