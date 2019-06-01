package webserver

import (
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Error Objects
var (
	errNotFound           = errors.New("not found")
	errUpdatedBoth        = errors.New("you can not update short and root link at once")
	errShortAlreadyExists = errors.New("the set short identifyer already exists")
	errInvalidArguments   = errors.New("invalid arguments")
)

// Static File Handlers
var (
	fileHandlerStatic = fasthttp.FS{
		Root:       "./web/dist",
		IndexNames: []string{"index.html"},
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			return ctx.Path()[7:]
		},
	}
)

type Config struct {
	Addr string     `json:"addr"`
	TLS  *TLSConfig `json:"tls"`
}

type TLSConfig struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"certfile"`
	Key     string `json:"keyfile"`
}

type WebServer struct {
	server *fasthttp.Server
	router *routing.Router
}
