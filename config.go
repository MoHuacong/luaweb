package luaweb

import (
	"fmt"
)

type Config struct {
	api *Api
}

func NewConfig(api *Api) (*Config, error) {
	return &Config{api}, nil
}

func (config *Config) GetData() string {
	fmt.Println(config.api.web.GetDir())
	fmt.Println(config.api.req.Host)
	fmt.Println(config.api.FormatFile(config.api.req.Host))
	return ""
}

func (config *Config) Analyze() {
	config.GetData()
}