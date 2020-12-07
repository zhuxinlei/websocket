package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/vrecan/death"
	"server-tokenhouse-ws/config"
	"server-tokenhouse-ws/service"
	"server-tokenhouse-ws/ws"
	"syscall"
)

func main() {
	err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	cfg := config.GetServer()

	setLogLevel(cfg.LogLevel)

	// 初始化主题校验
	ws.InitValidTopicTrie()

	go service.Run(cfg.Host, cfg.Port)
	hubManager := ws.GetHubManager()
	go hubManager.Run()

	log.Println("token house websocket started.")

	// 捕捉退出信号
	d := death.NewDeath(syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGALRM)
	d.WaitForDeathWithFunc(func() {
		log.Println("token house websocket stopped.")
	})
}

func setLogLevel(lvl string) {
	switch lvl {
	case "Debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
