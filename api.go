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

type Api struct {
	web *Web
	req *http.Request
	resp http.ResponseWriter
	key map[interface{}]interface{}
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func NewApi(web *Web, req *http.Request, resp http.ResponseWriter) (*Api, error) {

	api := &Api{web, req, resp, make(map[interface{}]interface{})}
	
	var port int
	var arr []string = strings.Split(api.req.Host, ":");
	
	if(len(arr) != 2) {
		port = 80
	} else {
		port, _ = strconv.Atoi(arr[1])
	}
	
	if api.web.GetApiPort() != port {
		return api, nil
	}
	
	api.inits()
	api.main()
	return api, nil
}

func (api *Api) Key(key string) interface{} {
	return api.key[key]
}

func (api *Api) main() error {
	json := make(map[string]interface{})
	types := api.req.URL.Query().Get("type")
	
	/* 登录 */
	if types == "login" {
		user := api.req.URL.Query().Get("user")
		pass := api.req.URL.Query().Get("pass")
		ret := api.Login(user, pass)
		if ret == true {
			json["code"] = 0
			json["data"] = "登录成功"
		} else {
			json["code"] = 1
			json["data"] = "登录失败"
		}
		return api.JsonExit(json)
	}
	
	/* 检测登录 */
	if !api.IsLoginCookie() {
		json["code"] = 1
		json["data"] = "未登陆"
		return api.JsonExit(json)
	}
	
	/* 端口监听 */
	if types == "AddListen" {
		api.AddListen(api.req.URL.Query().Get("host"))
		json["code"] = 0
		json["data"] = "添加成功"
		return api.JsonExit(json)
	}
	
	json["code"] = 1
	json["data"] = "作者Moid最帅"
	
	return api.JsonExit(json)
}

/* 初始并返回json字符串 */
func (api *Api) json() string {
	
	api.key["user"] = "admin"
	api.key["pass"] = "admin888"
	
	mjson, _ := json.Marshal(api.key)
	return string(mjson)
}

/* 初始化 */
func (api *Api) inits() {
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
func (api *Api) Login(user string, pass string) bool {
	if user != api.key["user"] || Md5(pass) != Md5(api.key["pass"].(string)) {
		return false
	}
	
	userc := http.Cookie{Name: "user", Value: user, MaxAge: 60 * 60 * 24 * 3600}
	passc := http.Cookie{Name: "pass", Value: Md5(pass), MaxAge: 60 * 60 * 24 * 3600}
	http.SetCookie(api.resp, &userc)
	http.SetCookie(api.resp, &passc)
	return true
}

/* 添加端口监听 */
func (api *Api) AddListen(host string) bool {
	var err error
	
	go func() {
		err = api.web.AddListen(host)
	}()
	
	if err != nil {
		return false
	}
	return true
}