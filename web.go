package luaweb

import (
	"os"
	"fmt"
	"strconv"
	"net/http"
	"strings"
	//lua "github.com/yuin/gopher-lua"
)

type Web struct {
	fd uint
	host string
	api_host string
	dir string
}

func NewWeb(host string, api_host string, dir string) (*Web, error) {
	return &Web{0, host, api_host, dir}, nil
}

func (web *Web) ServeHTTP(resp http.ResponseWriter,req *http.Request) {
	NewApi(web, req, resp)
	config, _ := NewConfig(web.dir, req, resp)
	config.Analyze()
}

func (web *Web) start() {
	fmt.Println("ok")
	_, err := os.Stat(web.dir)
	
	if os.IsNotExist(err) {
		err = os.Mkdir(web.dir, os.ModePerm)
	}
	
	go http.ListenAndServe(web.host, web)
	http.ListenAndServe(web.api_host, web)
}

func (web *Web) GetDir() string {
	return web.dir
}

func (web *Web) GetApiHost() string {
	return web.api_host
}

func (web *Web) GetApiIP() string {
	return strings.Split(web.GetApiHost(), ":")[0]
}

func (web *Web) GetApiPort() int {
	ret, _ := strconv.Atoi(strings.Split(web.GetApiHost(), ":")[1])
	return ret
}

func (web *Web) AddListen(host string) error {
	return http.ListenAndServe(host, web)
}

func (web *Web) NonBlocking() {
	go web.start()
}

func (web *Web) Run(status bool) {
	if status == true {
		web.start()
	} else {
		web.NonBlocking()
	}
}