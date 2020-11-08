package config

import (
	//"fmt"
	"io/ioutil"
	"github.com/flosch/pongo2"
	. "github.com/MoHuacong/luaweb/api"
)

type Config struct {
	Api *Api
}

func NewConfig(api *Api) (*Config, error) {
	return &Config{api}, nil
}

func (config *Config) GetData() string {
	path := config.Api.FormatFile(config.Api.Req.Host)
	ch, err := ioutil.ReadFile(path)
	if err != nil { return "" }
	return string(ch)
}

func (config *Config) Analyze() bool {
	user := config.GetData()
	if user == "" {
		config.MainService()
		return false
	} else {
		web_dir := config.Api.Key("web_dir").(string)
		config.Api.Value("web", web_dir + "/" + user)
		return true
	}
}

/* 引导页 */
func (config *Config) MainService() {
	tpl_dir := config.Api.Key("tpl_dir").(string)
	ch, err := ioutil.ReadFile(tpl_dir + "/index.html")
	
	html := "域名" + config.Api.Req.Host + "没有与服务绑定"
	
	var data string
	if err == nil {
		data = string(ch)
	} else {
		data = "域名{{host}}没有与服务绑定"
	}
	
	tpl, err1 := pongo2.FromString(data)
	
	var out string
	if err1 != nil {
		config.Api.Resp.Write([]byte(data))
	} else if out, err = tpl.Execute(pongo2.Context{"host": config.Api.Req.Host}); err != nil{
		config.Api.Resp.Write([]byte(html))
	} else {
		config.Api.Resp.Write([]byte(out))
	}
}