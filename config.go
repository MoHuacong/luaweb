package luaweb

import (
	"fmt"
	"io/ioutil"
	"github.com/flosch/pongo2"
)

type Config struct {
	api *Api
}

func NewConfig(api *Api) (*Config, error) {
	return &Config{api}, nil
}

func (config *Config) GetData() string {
	path := config.api.FormatFile(config.api.req.Host)
	ch, err := ioutil.ReadFile(path)
	if err != nil { return "" }
	return string(ch)
}

func (config *Config) Analyze() {
	data := config.GetData()
	if data == "" {
		config.MainService()
	} else {
		fmt.Println("ok")
	}
}

func (config *Config) MainService() {
	tpl_dir := config.api.Key("tpl_dir").(string)
	ch, err := ioutil.ReadFile(tpl_dir + "/index.html")
	
	html := "域名" + config.api.req.Host + "没有与服务绑定"
	
	var data string
	if err == nil {
		data = string(ch)
	} else {
		data = "域名{{host}}没有与服务绑定"
	}
	
	tpl, err1 := pongo2.FromString(data)
	
	var out string
	if err1 != nil {
		config.api.resp.Write([]byte(data))
	} else if out, err = tpl.Execute(pongo2.Context{"host": config.api.req.Host}); err != nil{
		config.api.resp.Write([]byte(html))
	} else {
		config.api.resp.Write([]byte(out))
	}
}