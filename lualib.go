package luaweb

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	lib "github.com/MoHuacong/luaweb/lualib"
)

type Lualib struct {
	L *lua.LState
	config *Config
}

func NewLualib(config *Config) (*Lualib, error) {
	this := &Lualib{lua.NewState(), config}
	this.InitFunc()
	return this, nil
}

func (this *Lualib) InitFunc() {
	lib.HTML.Data = ""
	this.L.SetGlobal("double", this.L.NewFunction(lib.Print))
}

func (this *Lualib) Main() {
	fmt.Println(this.config.api.Key("web"))
}

func (this *Lualib) Close() {
	this.L.Close()
}