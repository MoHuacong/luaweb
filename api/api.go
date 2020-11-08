package api

import (
	"os"
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"io/ioutil"
	"net/url"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

/* 配置数据结构体 */
type ConfigData struct {
	File_End string
	Host_dir string
	FileName string
}

var CData ConfigData = ConfigData{
			File_End: ".json",
			Host_dir: "host",
			FileName: "data",
			}

type Api struct {
	A bool
	ret string
	data map[string]interface{}
	Req *http.Request
	err map[int]string
	Resp http.ResponseWriter
	key map[interface{}]interface{}
}

/* md5加密 */
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/* 创建并初始化结构体 */
func NewApi(data map[string]interface{}, req *http.Request, resp http.ResponseWriter) (*Api, error) {

	api := &Api{
	A: false,
	ret: "",
	data: data,
	Req: req,
	Resp: resp,
	err: make(map[int]string),
	key: make(map[interface{}]interface{}),
	}
	
	api.inits()
	
	var port int
	var arr []string = strings.Split(api.Req.Host, ":")
	
	if(len(arr) != 2) {
		port = 80
	} else {
		port, _ = strconv.Atoi(arr[1])
	}
	
	if api.GetApiPort() != port {
		return api, nil
	}
	
	api.A = true
	
	api.main()
	return api, nil
}

/* value获取 */
func (api *Api) Key(key string) interface{} {
	return api.key[key]
}

/* key/value设置 */
func (api *Api) Value(key string, value interface{}) bool {
	api.key[key] = value
	if api.key[key] != value {
		return false
	}
	return true
}

/* 强制类型转换 */
func (api *Api) ToType(mt reflect.Type, url url.Values) []reflect.Value {
	var args []reflect.Value = make([]reflect.Value, mt.NumIn())
	
	var i int = 0
	for key, value := range url {
		if key == "type" { continue }
		if i >= mt.NumIn() { break }
		
		kind := mt.In(i).Kind()
		t := kind.String()

		if t == "string" {
			args[i] = reflect.ValueOf(value[0])
		}
		i++
	}
	return args
}

/* 反射方法 */
func (api *Api) reflectMethod(method string) uint {
	defer func() uint {
		if p := recover(); p != nil {
			fmt.Printf("panic错误: %s\n", p)
			return api.errExit(1, "panic错误: " + p.(string))
		}
		return api.retPut(api.ret)
	}()
	
	v := reflect.ValueOf(api)
	if method == "" {
		return api.errExit(3, "type参数不完整")
	}
	
	var zero reflect.Value
	mv := v.MethodByName(method)
	if mv == zero {
		return api.errExit(3, "type接口错误或不存在(反射失败)")
	}
	
	mt := mv.Type()
	if mt.NumIn() > len(api.Req.URL.Query()) - 1 {
		return api.errExit(2, "参数不完整")
	}
	
	args := api.ToType(mt, api.Req.URL.Query())

	value := mv.Call(args)
	return uint(value[0].Uint())
}

/* 路由器 */
func (api *Api) main() error {
	types := api.Req.URL.Query().Get("type")
	
	if types != "Login" {
		/* 检测登录 */
		if !api.IsLoginCookie() {
			return api.FormatMapExit(1, "末登录")
		}
	}
	
	ret := api.reflectMethod(types)
	if ret  == 0 {
		return api.FormatMapExit(ret, api.ret)
	} else {
		return api.FormatMapExit(ret, api.Errors())
	}
	
	return api.FormatMapExit(1, "作者Moid最帅")
}

/* 初始并返回json字符串 */
func (api *Api) json() string {
	
	api.key["user"] = "admin"
	api.key["pass"] = "admin888"
	
	key := make(map[string]interface{})
	key["user"] = api.key["user"].(string)
	key["pass"] = api.key["pass"].(string)
	mjson, _ := json.Marshal(key)
	return string(mjson)
}

/* 初始化目录 */
func (api *Api) initDir() bool {
	dir := api.GetDir()
	
	tpl := dir + "/tpl"
	host := dir + "/host"
	user := dir + "/user"
	web := dir + "/wwwroot"
	
	os.Mkdir(tpl, os.ModePerm)
	os.Mkdir(host, os.ModePerm)
	os.Mkdir(user, os.ModePerm)
	os.Mkdir(web, os.ModePerm)
	
	api.key["tpl_dir"] = tpl
	api.key["host_dir"] = host
	api.key["user_dir"] = user
	api.key["web_dir"] = web
	
	return true
}

/* 初始化 */
func (api *Api) inits() {
	api.initDir()
	var path string = api.GetDir() + "/data.json"
	_, err := os.Stat(path)
	
	if os.IsNotExist(err) {
		f, _ := os.Create(path)
		f.WriteString(api.json())
		f.Close()
	} else {
		key := make(map[string]interface{})
		ch, _ := ioutil.ReadFile(path)
		json.Unmarshal(ch, &key)
		
		for k, v := range key {
			api.key[k] = v.(string)
		}
	}
}

/* 直接格式化返回 */
func (api *Api) FormatMapExit(code uint, data interface{}) error {
	json := make(map[string]interface{})
	json["code"] = code
	json["data"] = data
	return api.JsonExit(json)
	
}

/* json格式返回 */
func (api *Api) JsonExit(v interface{}) error {
	mjson, _ := json.Marshal(v)
	_, err := api.Resp.Write(mjson)
	return err
}

/* 登录判断 */
func (api *Api) IsLoginCookie() bool {
	user, err := api.Req.Cookie("user")
	if err != nil { return false }
	pass, err2 := api.Req.Cookie("pass")
	if err2 != nil { return false }
	if user.Value != api.key["user"] || pass.Value != Md5(api.key["pass"].(string)) {
		return false
	}
	return true
}

/* 登录 */
func (api *Api) Login(user string, pass string) uint {
	if user != api.key["user"] || Md5(pass) != Md5(api.key["pass"].(string)) {
		return api.errExit(1, "帐号或密码不正确")
	}
	
	userc := http.Cookie{Name: "user", Value: user, MaxAge: 60 * 60 * 24 * 3600}
	passc := http.Cookie{Name: "pass", Value: Md5(pass), MaxAge: 60 * 60 * 24 * 3600}
	http.SetCookie(api.Resp, &userc)
	http.SetCookie(api.Resp, &passc)
	return api.retPut("登录成功")
}

/* 添加用户 */
func (api *Api) AddUser(user string, pass string) uint {
	file := api.Key("user_dir").(string) + "/" + user + CData.File_End
	ch, err := ioutil.ReadFile(file)
	key := make(map[string]interface{})
	
	if err == nil {
		json.Unmarshal(ch, &key)
		if key["user"] != nil {
			return api.errExit(1, "已有[" + key["user"].(string) + "]用户")
		}
	}
	
	key = make(map[string]interface{})
	key["user"] = user
	key["pass"] = Md5(pass)
	
	mjson, err2 := json.Marshal(key)
	
	if err2 != nil {
		return api.errExit(2, "生成用户失败")
	}
	
	f, err3 := os.Create(file)
	if err3 != nil {
		return api.errExit(3, "添加用户失败")
	}
	f.WriteString(string(mjson))
	f.Close()
	
	dir := api.Key("web_dir").(string) + "/" + user
	err = os.Mkdir(dir, os.ModePerm)
	
	if err != nil {
		return api.errExit(4, "创建目录失败")
	}
	
	return api.retPut("添加用户成功")
}

/* 添加用户域名 */
func (api *Api) AddUserHost(user string, host string) uint {
	ch, err := ioutil.ReadFile(api.FormatFile(host))
	
	if err == nil && string(ch) != user {
		return api.errExit(1, "域名已绑定")
	}
	
	file := api.Key("user_dir").(string) + "/" + user + CData.File_End
	ch, err = ioutil.ReadFile(file)
	
	if err != nil {
		return api.errExit(2, "用户不存在")
	}
	
	key := make(map[string]interface{})
	err = json.Unmarshal(ch, &key)
	if err != nil {
		return api.errExit(3, "用户数据错误")
	}
	
	if key["host"] != nil {
		for _, v := range key["host"].([]interface{}) {
			if v.(string) == host {
				return api.errExit(4, "域名已添加")
			}
		}
	}
	
	if key["host"] == nil {
		key["host"] = make([]interface{}, 1)
		key["host"].([]interface{})[0] = host
	} else {
		arr := key["host"].([]interface{})
		key["host"] = append(arr, host)
	}
	
	if api.PutJsonData(file, key) != 0 {
		return api.errExit(5, "用户主机数据写入失败")
	}
	
	if api.PutData(api.FormatFile(host), user) != 0 {
		return api.errExit(6, "添加用户主机失败")
	}
	
	return api.retPut("添加成功")
}

/* 添加端口监听 */
/*
func (api *Api) AddListen(host string) uint {
	var err error
	
	go func() {
		err = api.web.AddListen(host)
	}()
	
	if err != nil {
		return api.errExit(1, "监听失败")
	}
	return api.retPut("监听成功")
}
*/

func (api *Api) GetDir() string {
	return api.data["dir"].(string)
}

func (api *Api) GetApiHost() string {
	return api.data["server"].(string)
}

func (api *Api) GetApiIP() string {
	return strings.Split(api.GetApiHost(), ":")[0]
}

func (api *Api) GetApiPort() int {
	ret, _ := strconv.Atoi(strings.Split(api.GetApiHost(), ":")[1])
	return ret
}

/* 写入数据 */
func (api *Api) PutJsonData(file string, data map[string]interface{}) uint {
	mjson, err := json.Marshal(data)
	
	if err != nil {
		return api.errExit(1, "生成失败")
	}
	
	if api.PutData(file, string(mjson)) != 0 {
		return api.errExit(2, "写入失败")
	}
	
	return api.retPut("写入数据成功")
}

/* 写入raw数据 */
func (api *Api) PutData(file string, data string) uint {
	
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		var dir string
		arr := strings.Split(file, "/")
		for k, v := range arr {
			if k + 1 != len(arr) {
				dir += v + "/"
				os.Mkdir(dir, os.ModePerm)
			}
		}
		f, err1 := os.Create(file)
		if err1 != nil {
			return api.errExit(1, "写入失败")
		}
		_, err = f.WriteString(data)
		if err != nil {
			return api.errExit(2, "写入失败")
		}
		f.Close()
	} else {
		f, _ := os.OpenFile(file, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, os.ModePerm)
		f.WriteString(data)
		f.Close()
	}
	
	return api.retPut("写入数据成功")
}

/* 域名转路径 */
func (api *Api) FormatDir(host string) string {
	var port string
	arr := strings.Split(host, ":")
	if len(arr) == 2 {
		port = arr[1]
		arr = strings.Split(arr[0], ".")
	} else {
		port = "80"
		arr = strings.Split(host, ".")
	}
	
	var dir string = port + "/"
	for i := len(arr) - 1; i >= 0; i-- {
		dir += arr[i] + "/"
	}
	return dir
}

/* 域名转绝对路径 */
func (api *Api) FormatDirs(host string) string {
	return api.data["dir"].(string) + "/" + CData.Host_dir + "/" + api.FormatDir(host)
}

/* 域名转绝对文件路径 */
func (api *Api) FormatFile(host string) string {
	return api.FormatDirs(host) + CData.FileName + CData.File_End
}

/* 返回值写入 */
func (api *Api) retPut(str string) uint {
	api.ret = str
	return 0
}

/* 直接写入并返回原值 */
func (api *Api) errExit(code uint, data string) uint {
	api.errPut(data)
	return code
}
/* 写入错误信息 */
func (api *Api) errPut(str string) bool {
	api.err[len(api.err)] = str
	if api.err[len(api.err) - 1] != str {
		return false
	}
	return true
}

/* 输出错误信息 */
func (api *Api) Errors() string {
	l := len(api.err)
	if l <= 0 {
		return ""
	}
	return api.err[l - 1]
}