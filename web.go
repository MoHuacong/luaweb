package luaweb

import (
	"fmt"
	. "github.com/MoHuacong/luaweb/api"
	. "github.com/MoHuacong/luaweb/config"
	. "github.com/MoHuacong/luaweb/lualib"
	"net/http"
	"os"
)

type Web struct {
	data map[string]interface{}
}

func NewWeb(data map[string]interface{}) (*Web, error) {
	return &Web{data}, nil
}

func (web *Web) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	api, _ := NewApi(web.data, req, resp)
	if !api.A {
		config, _ := NewConfig(api)
		if config.Analyze() {
			lib, _ := NewLualib(config)
			lib.Main()
			lib.Close()
		}
	}
}

func (web *Web) start() {
	fmt.Println("ok")
	_, err := os.Stat(web.data["dir"].(string))

	if os.IsNotExist(err) {
		err = os.Mkdir(web.data["dir"].(string), os.ModePerm)
	}

	go http.ListenAndServe(web.data["host"].(string), web)
	http.ListenAndServe(web.data["server"].(string), web)
}

func (web *Web) NonBlocking() {
	go web.start()
}

func (web *Web) AddListen(host string) error {
	return http.ListenAndServe(host, web)
}

func (web *Web) Run(status bool) {
	if status == true {
		web.start()
	} else {
		web.NonBlocking()
	}
}
