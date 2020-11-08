package lualib

import (
	. "github.com/yuin/gopher-lua"
)

func RegisterModule(ls *LState, name string, funcs map[string]LGFunction) LValue {
	tb := ls.FindTable(ls.Get(RegistryIndex).(*LTable), "_LOADED", 1)
	newmod := ls.FindTable(ls.Get(GlobalsIndex).(*LTable), name, len(funcs))
	newmodtb, _ := newmod.(*LTable);
	for fname, fn := range funcs {
		newmodtb.RawSetString(fname, ls.NewFunction(fn))
	}
	ls.SetField(tb, name, newmodtb)
	return newmodtb
}