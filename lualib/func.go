package lualib

import (
	"fmt"
	"strings"
	. "github.com/yuin/gopher-lua"
)

func Print(L *LState) int {
	top := L.GetTop()
	resp := L.Options.LuaWeb.Resp
	for i := 1; i <= top; i++ {
		resp.Write([]byte(L.ToStringMeta(L.Get(i)).String()))
		resp.Write([]byte("</br>"))
	}
	return 0
}

func Printf(L *LState) int {
	top := L.GetTop()
	resp := L.Options.LuaWeb.Resp
	for i := 1; i <= top; i++ {
		resp.Write([]byte(L.ToStringMeta(L.Get(i)).String()))
	}
	return 0
}

func Sprintf(L *LState) int {
	str := L.CheckString(1)
	resp := L.Options.LuaWeb.Resp
	args := make([]interface{}, L.GetTop()-1)
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		args[i-2] = L.Get(i)
	}
	npat := strings.Count(str, "%") - strings.Count(str, "%%")
	resp.Write([]byte(fmt.Sprintf(str, args[:intMin(npat, len(args))]...)))
	return 0
}

func Init(L *LState) {
	OpenLibs(L)
	L.SetGlobal("print", L.NewFunction(Print))
	L.SetGlobal("printf", L.NewFunction(Printf))
	L.SetGlobal("sprintf", L.NewFunction(Sprintf))
}