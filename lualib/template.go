package lualib

import (
	"os"
	"fmt"
	"reflect"
	"unsafe"
	"github.com/flosch/pongo2"
	. "github.com/yuin/gopher-lua"
	. "github.com/MoHuacong/luaweb/context"
)

type Template struct {
	Cxt *Context
	Data pongo2.Context
	Tpl *pongo2.Template
}

const luaTemplateTypeName = "template"

var templateFuncs = map[string]LGFunction{
	"test":     templateTest,
	"new":     templateNew,
}

var templateMtFuncs = map[string]LGFunction{
	"set": mtSet,
	"out": mtOut,
}

func inits(L *LState, filename string) *Template {
	var err error
	var tpl *pongo2.Template
	_, err = os.Open(filename)
	if err != nil && os.IsNotExist(err) {
		tpl, err = pongo2.FromString(filename)
	} else {
		tpl, err = pongo2.FromFile(filename)
	}
	
	if err != nil {
		L.ArgError(1, "new template expected")
	}
	return &Template{
		Cxt: L.Options.LuaWeb,
		Data: make(pongo2.Context),
		Tpl: tpl,
	};
}

func checkTemplate(L *LState) *Template {
	ud := L.CheckUserData(1)
	if tp, ok := ud.Value.(*Template); ok {
		return tp
	}
	L.ArgError(1, "template expected")
	return nil
}

func TableValue(tb *LTable, name string) reflect.Value {
	v := reflect.ValueOf(tb)
	if v.Kind() != reflect.Ptr {
		var lnil reflect.Value
		return lnil
	}
	
	v = v.Elem()
	return v.FieldByName(name)
	
	/*
	if name.Kind() == reflect.String {
		name.SetString("小学生")
	}
	*/
}

func TableArray(tb *LTable, name string) []LValue {
	var data []LValue
	n := TableValue(tb, name)
	if n.Kind() != reflect.Slice && n.Kind() != reflect.Array {
		return data
	}
	ptr := n.Pointer()
	var v LValue
	
	for i := 0; i < n.Len(); i++ {
		v = *(*LValue)(unsafe.Pointer(ptr))
		data = append(data, v)
		ptr += unsafe.Sizeof(v)
	}
	
	return data
}

func TableMap(tb *LTable, name string) map[LValue]LValue {
	n := TableValue(tb, name)
	data := make(map[LValue]LValue, 1)
	if n.Kind() != reflect.Map {
		return data
	}
	var key LValue
	if n.MapKeys() == nil {
		return data
	}
	for _, k := range n.MapKeys() {
		v := n.MapIndex(k)
		switch k.Kind() {
			case reflect.Int, reflect.Uint:
				key = LNumber(k.Int())
			case reflect.String:
				key = LString(k.String())
			default:
		}
		p := v.InterfaceData()
		data[key] = *(*LValue)(unsafe.Pointer(&p))
	}
	return data
}

func Analyze(tb *LTable) map[string]interface{} {
	data := make(map[string]interface{})
	arr := TableArray(tb, "array")
	strdict := TableMap(tb, "strdict")
	dict := TableMap(tb, "dict")
	for k, v := range arr {
		if v.Type() == LTTable {
			data[LNumber(k).String()] = Analyze(v.(*LTable))
		} else {
			data[LNumber(k).String()] = v
		}
	}
	for k, v := range strdict {
		if v.Type() == LTTable {
			data[k.String()] = Analyze(v.(*LTable))
		} else {
			data[k.String()] = v
		}
	}
	for k, v := range dict {
		if v.Type() == LTTable {
			data[k.String()] = Analyze(v.(*LTable))
		} else {
			data[k.String()] = v
		}
	}
	return data
}

func OpenTemplate(L *LState) int {
	tplmod := RegisterModule(L, TemplateLibName, templateFuncs)
	L.Push(tplmod)
	return 1
}

func templateTest(L *LState) int {
	tb := L.CheckTable(1)
	fmt.Println(Analyze(tb))
	L.Push(LString("Moid"))
	return 1
}

func templateNew(L *LState) int {
	filename := L.ToString(1)
	ud := L.NewUserData()
	ud.Value = inits(L, filename)
	mt := L.NewTypeMetatable(luaTemplateTypeName)
	mt.RawSetString("__index", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), templateMtFuncs))
	L.SetMetatable(ud, mt)
	L.Push(ud)
	return 1
}

func mtSet(L *LState) int {
	tpl := checkTemplate(L)
	name := L.ToString(2)
	value := L.Get(3)
	
	if value.Type() == LTTable && name != "" {
		tpl.Data[name] = Analyze(value.(*LTable))
		L.Push(LBool(true))
		return 1
	}
	
	if name != ""{
		tpl.Data[name] = value.String()
		L.Push(LBool(true))
	} else {
		L.Push(LBool(false))
	}
	return 1
}

func mtOut(L *LState) int {
	tpl := checkTemplate(L)
	out, err := tpl.Tpl.Execute(tpl.Data)
	if err != nil {
		L.Push(LBool(false))
	}
	L.Push(LString(out))
	return 1
}