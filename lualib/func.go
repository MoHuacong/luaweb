package lualib

import (
	//"fmt"
	. "github.com/yuin/gopher-lua"
)

type IO struct {
	Data string
}

var HTML IO = IO{""}

func Print(L *LState) int {
	top := L.GetTop()
	for i := 1; i <= top; i++ {
		HTML.Data += L.ToStringMeta(L.Get(i)).String()
		if i != top {
			HTML.Data += "</br>"
		}
	}
	return 0
}