package lualib

import (
	"os"
	//"fmt"
	"strings"
	"io/ioutil"
	lua "github.com/yuin/gopher-lua"
	"github.com/MoHuacong/luaweb/context"
	. "github.com/MoHuacong/luaweb/config"
)

type Lualib struct {
	L *lua.LState
	config *Config
}

func NewLualib(config *Config) (*Lualib, error) {
	InitType()
	config.Api.Resp.Header().Add("Content-Type", "text/html");
	config.Api.Req.ParseMultipartForm(1024*1024)
	cxt, _ := context.NewContext(config.Api.Key("web").(string), config, config.Api.Req, config.Api.Resp)
	op := lua.Options{LuaWeb: cxt}
	this := &Lualib{lua.NewState(op), config}
	this.Init()
	return this, nil
}

func (this *Lualib) Init() {
	Init(this.L)
}

func (this *Lualib) Main() bool {
	os.Mkdir(this.config.Api.Key("web").(string), os.ModePerm)
	file := this.config.Api.Key("web").(string) + this.config.Api.Req.URL.Path
	_, err := os.Open(file)
	if err != nil && os.IsNotExist(err) {
		this.L.DoString(`print("<h1>404</h1>")`)
		return true
	}
	
	arr := strings.Split(file, ".")
	end := arr[len(arr)-1]
	
	this.config.Api.Resp.Header().Add("Content-Type", HttpType.Type(end));
	
	if end == "lua" {
		err = this.L.DoFile(file)
		if err != nil {
			this.config.Api.Resp.Write([]byte(err.Error()))
			return false
		}
		return true
	}
	ch, _ := ioutil.ReadFile(file)
	this.config.Api.Resp.Write(ch)
	return true
}

func (this *Lualib) Close() {
	this.L.Close()
}