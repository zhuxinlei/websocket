package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var srvManger *Manger

type Manger struct {
	server *Server
}

func GetServer() *Server {
	return srvManger.server
}

// 加载配置文件
func LoadConfig(path string) error {
	// 读取基本配置
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	serve := Server{}
	err = yaml.Unmarshal(data, &serve)
	if err != nil {
		return err
	}

	srvManger = &Manger{
		server: &serve,
	}

	return nil
}
