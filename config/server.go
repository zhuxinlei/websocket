package config

// 服务配置
type Server struct {
	Host           string `yaml:"host"` // 主机地址
	Port           int    `yaml:"port"` // 端口号
	LogLevel       string `yaml:"logLevel"`
	Topic          Topic  `yaml:"topic"`          // 主题配置
	TokenHouseHost string `yaml:"tokenHouseHost"` //token
}

type Topic struct {
	Valid []string            `yaml:"valid"`
	Dict  map[string][]string `yaml:"dict"`
}
