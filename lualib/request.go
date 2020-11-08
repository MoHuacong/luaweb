package lualib

import (
	"io"
	"os"
	//"fmt"
	"net/url"
	"strings"
	"mime/multipart"
	"layeh.com/gopher-luar"
	. "github.com/yuin/gopher-lua"
)

var requestFuncs = map[string]LGFunction{
	"file":     reqFile,
	"get":     reqGet,
	"post":    reqPost,
	"cookie":  reqCookie,
	"header":  reqHeader,
}

func valueToTable(value url.Values) LTable {
	var data LTable
	for k, v := range value {
		data.RawSet(LString(k), LString(v[0]))
	}
	return data
}

func OpenRequest(L *LState) int {
	reqmod := RegisterModule(L, RequestLibName, requestFuncs)
	L.Push(reqmod)
	return 1
}

func reqGet(L *LState) int {
	key := L.ToString(1)
	req := L.Options.LuaWeb.Req
	values := req.URL.Query()
	if key != "" {
		v := values.Get(key)
		if v != "" {
			L.Push(LString(v))
			return 1
		}
		L.Push(LNil)
		return 1
	}
	table := valueToTable(values)
	L.Push(&table)
	return 1
}

func reqPost(L *LState) int {
	key := L.ToString(1)
	req := L.Options.LuaWeb.Req
	if key != "" {
		v := req.PostFormValue(key)
		
		if v != "" {
			L.Push(LString(v))
			return 1
		}
		L.Push(LNil)
		return 1
	}
	table := valueToTable(req.PostForm)
	L.Push(&table)
	return 1
}

func reqCookie(L *LState) int {
	var data LTable
	key := L.ToString(1)
	req := L.Options.LuaWeb.Req
	if key != "" {
		cookie, _ := req.Cookie(key)
		L.Push(luar.New(L, cookie))
		return 1
	}
	for _, v := range req.Cookies() {
		data.RawSet(LString(v.Name), luar.New(L, v))
	}
	L.Push(&data)
	return 1
}

var fileFuncs = map[string]LGFunction{
	"size":      fileSize,
	"type":      fileType,
	"save":     fileSave,
	"name":    fileName,
}

func checkReqFile(L *LState) *multipart.FileHeader {
	ud := L.CheckUserData(1)
	if head, ok := ud.Value.(*multipart.FileHeader); ok {
		return head
	}
	L.ArgError(1, "template expected")
	return nil
}

func newFile(L *LState, filehead *multipart.FileHeader) *LUserData {
	ud := L.NewUserData()
	ud.Value = filehead
	mt := L.NewTypeMetatable("filehead")
	mt.RawSetString("__index", mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), fileFuncs))
	L.SetMetatable(ud, mt)
	return ud
}

func fileName(L *LState) int {
	filehead := checkReqFile(L)
	L.Push(LString(filehead.Filename))
	return 1
}

func fileSize(L *LState) int {
	filehead := checkReqFile(L)
	L.Push(LNumber(filehead.Size))
	return 1
}

func fileType(L *LState) int {
	var end string
	filehead := checkReqFile(L)
	arr := strings.Split(filehead.Filename, ".")
	if arr == nil {
		end = filehead.Filename
	} else {
		end = arr[len(arr)-1]
	}
	L.Push(LString(HttpType.Type(end)))
	return 1
}

func fileSave(L *LState) int {
	filehead := checkReqFile(L)
	filename := L.ToString(2)
	
	var err error
	var f *os.File
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
	} else {
		f, err = os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, os.ModePerm)
	}
	defer f.Close()
	file, _ := filehead.Open()
	_, err = io.Copy(f, file)
	if err != nil {
		L.Push(LBool(false))
		return 1
	}
	L.Push(LBool(true))
	return 1
}

func reqFile(L *LState) int {
	var data LTable
	key := L.ToString(1)
	req := L.Options.LuaWeb.Req
	if key != "" {
		_, head, err := req.FormFile(key)
		if head != nil && err == nil {
			L.Push(newFile(L, head))
			return 1
		}
		L.Push(LNil)
		return 1
	}
	for name, arr := range req.MultipartForm.File {
		var tb LTable
		for k, filehead := range arr {
			tb.RawSetInt(k+1, newFile(L, filehead))
		}
		data.RawSet(LString(name), &tb)
	}
	L.Push(&data)
	return 1
}

func reqHeader(L *LState) int {
	name := L.ToString(1)
	req := L.Options.LuaWeb.Req
	value := req.Header.Get(name)
	if name == "" || value == "" {
		L.Push(LNil)
		return 1
	}
	L.Push(LString(value))
	return 1
}