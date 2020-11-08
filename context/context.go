package context

import (
	"net/http"
	. "github.com/MoHuacong/luaweb/config"
)

type Context struct {
	Dir string
	Config *Config
	Req *http.Request
	Resp http.ResponseWriter
}

func NewContext(dir string, config *Config, req *http.Request, resp http.ResponseWriter) (*Context, error) {
	return &Context{
		Dir: dir,
		Config: config,
		Req: req,
		Resp: resp,
		}, nil
}