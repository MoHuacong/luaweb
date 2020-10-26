package luaweb

import (
	"fmt"
	"net/http"
)

type Config struct {
	path string
	req *http.Request
	resp http.ResponseWriter
}

func NewConfig(path string, req *http.Request, resp http.ResponseWriter) (*Config, error) {
	return &Config{path, req, resp}, nil
}

func (config *Config) GetData() string {
	fmt.Println(config.path)
	fmt.Println(config.req.Host)
	return ""
}

func (config *Config) Analyze() {
	config.GetData()
}