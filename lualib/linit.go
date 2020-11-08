package lualib

import (
	//"fmt"
	. "github.com/yuin/gopher-lua"
)

const (
	LuaWebBaseLibName = ""
	LuaWebOsLibName = "os"
	RequestLibName = "request"
	TemplateLibName = "template"
)

type luaLib struct {
	libName string
	libFunc LGFunction
}

var luaLibs = []luaLib{
	luaLib{RequestLibName, OpenRequest},
	luaLib{TemplateLibName, OpenTemplate},
	luaLib{LuaWebOsLibName, LuaWebOpenOs},
}

func OpenLibs(ls *LState) {
	for _, lib := range luaLibs {
		//ls.PreloadModule(LuaWebOsLibName, lib.libFunc)
		ls.Push(ls.NewFunction(lib.libFunc))
		ls.Push(LString(lib.libName))
		ls.Call(1, 0)
	}
}
