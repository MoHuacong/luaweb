package luaweb

import (
	"os"
	//"fmt"
	"strings"
	"strconv"
	"io/ioutil"
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
	web *Web
	req *http.Request
	err map[int]string
	resp http.ResponseWriter
	key map[interface{}]interface{}
}

/* md5加密 */
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func NewApi(web *Web, req *http.Request, resp http.ResponseWriter) (*Api, error) {

	api := &Api{
	A: false,
	ret: "",
	web: web,
	req: req,
	resp: resp,
	err: make(map[int]string),
	key: make(map[interface{}]interface{}),
	}
	
	api.inits()
	
	var port int
	var arr []string = strings.Split(api.req.Host, ":")
	
	if(len(arr) != 2) {
		port = 80
	} else {
		port, _ = strconv.Atoi(arr[1])
	}
	
	if api.web.GetApiPort() != port {
		return api, nil
	}
	
	api.A = true
	
	api.main()
	return api, nil
}

func (api *Api) Key(key string) interface{} {
	return api.key[key]
}

func (api *Api) Value(key string, value interface{}) bool {
	api.key[key] = value
	if api.key[key] != value {
		return false
	}
	return true
}

/* 路由器 */
func (api *Api) main() error {
	var ret uint
	types := api.req.URL.Query().Get("type")
	
	/* 登录 */
	if types == "login" {
		user := api.req.URL.Query().Get("user")
		pass := api.req.URL.Query().Get("pass")
		ret = api.Login(user, pass)
		if ret  == 0 {
			return api.FormatMapExit(ret, api.ret)
		} else {
			return api.FormatMapExit(ret, api.Errors())
		}
	}
	
	/* 检测登录 */
	if !api.IsLoginCookie() {
		return api.FormatMapExit(1, "末登录")
	}
	
	/* 端口监听 */
	if types == "AddListen" {
		ret = api.AddListen(api.req.URL.Query().Get("host"))
		if ret  == 0 {
			return api.FormatMapExit(ret, api.ret)
		} else {
			return api.FormatMapExit(ret, api.Errors())
		}
	}
	
	/* 添加用户 */
	if types == "AddUser" {
		user := api.req.URL.Query().Get("user")
		pass := api.req.URL.Query().Get("pass")
		ret = api.AddUser(user, pass)
		if ret  == 0 {
			return api.FormatMapExit(ret, api.ret)
		} else {
			return api.FormatMapExit(ret, api.Errors())
		}
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
	dir := api.web.GetDir()
	
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
	var path string = api.web.GetDir() + "/data.json"
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
	_, err := api.resp.Write(mjson)
	return err
}

/* 登录判断 */
func (api *Api) IsLoginCookie() bool {
	user, err := api.req.Cookie("user")
	if err != nil { return false }
	pass, err2 := api.req.Cookie("pass")
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
	http.SetCookie(api.resp, &userc)
	http.SetCookie(api.resp, &passc)
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
	
	return api.retPut("添加用户成功")
}

/* 添加端口监听 */
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
	return api.web.dir + "/" + CData.Host_dir + "/" + api.FormatDir(host)
}

/* 域名转绝对文件路径 */
func (api *Api) FormatFile(host string) string {
	return api.FormatDirs(host) + "/" + CData.FileName + CData.File_End
}

/* 返回值写入 */
func (api *Api) retPut(str string) uint {
	api.ret = str
	return 0
}

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