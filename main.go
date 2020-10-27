// QUIC web server with built-in support for Lua, Markdown, Pongo2 and JSX.
package main

import (
	"fmt"
	"github.com/MoHuacong/luaweb"
)

func main() {
    web, _ := luaweb.NewWeb("0.0.0.0:8080", "0.0.0.0:2333", "/sdcard/luaweb")
    web.Run(true)
    fmt.Println(web)
}