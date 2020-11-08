package main

import (
	"fmt"
	"github.com/MoHuacong/luaweb"
)

func main() {
    arr := make(map[string]interface{})
    arr["host"] = "0.0.0.0:8080"
    arr["server"] = "0.0.0.0:2333"
    arr["dir"] = "/sdcard/luaweb"
    web, _ := luaweb.NewWeb(arr)
    web.Run(true)
    fmt.Println(web)
}