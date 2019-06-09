package webserver

import (
	"errors"

	"github.com/zekroTJA/myrunes/internal/database"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Error Objects
var (
	errNotFound         = errors.New("not found")
	errInvalidArguments = errors.New("invalid arguments")
	errUNameInUse       = errors.New("user name already in use")
	errNoAccess         = errors.New("access denied")
)

// Static File Handlers
var (
	fileHandlerStatic = fasthttp.FS{
		Root:       "./web/dist",
		IndexNames: []string{"index.html"},
		// PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
		// 	return ctx.Path()[7:]
		// },
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

	db   database.Middleware
	auth *Authorization

	config *Config
}

func NewWebServer(db database.Middleware, config *Config) (ws *WebServer) {
	ws = new(WebServer)

	ws.config = config
	ws.db = db
	ws.auth = NewAuthorization(db)
	ws.router = routing.New()
	ws.server = &fasthttp.Server{
		Handler: ws.router.HandleRequest,
	}

	ws.registerHandlers()

	return
}

func (ws *WebServer) registerHandlers() {
	ws.router.Use(ws.handlerFiles, ws.addCORSHeaders)

	api := ws.router.Group("/api")
	api.
		Post("/login", ws.handlerLogin)
	api.
		Post("/logout", ws.auth.LogOut)

	resources := api.Group("/resources")
	resources.
		Get("/champions", ws.handlerGetChamps)
	resources.
		Get("/runes", ws.handlerGetRunes)

	users := api.Group("/users")
	users.
		Post("", ws.handlerCreateUser)
	users.
		Post("/me", ws.auth.CheckRequestAuth, ws.handlerPostMe).
		Get(ws.auth.CheckRequestAuth, ws.handlerGetMe).
		Delete(ws.auth.CheckRequestAuth, ws.handlerDeleteMe)
	users.
		Get("/<uname>", ws.handlerCheckUsername)

	pages := api.Group("/pages", ws.addCORSHeaders, ws.auth.CheckRequestAuth)
	pages.
		Post("", ws.handlerCreatePage).
		Get(ws.handlerGetPages)
	pages.
		Get(`/<uid:\d+>`, ws.handlerGetPage).
		Post(ws.handlerEditPage).
		Delete(ws.handlerDeletePage)

	sessions := api.Group("/sessions", ws.addCORSHeaders, ws.auth.CheckRequestAuth)
	sessions.
		Get("", ws.handlerGetSessions)
	sessions.
		Delete(`/<id:\d+>`, ws.handlerDeleteSession)

	favorites := api.Group("/favorites", ws.addCORSHeaders, ws.auth.CheckRequestAuth)
	favorites.
		Get("", ws.handlerGetFavorites).
		Post(ws.handlerPostFavorite)
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.config.TLS

	if tls.Enabled {
		if tls.Cert == "" || tls.Key == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(ws.config.Addr, tls.Cert, tls.Key)
	}

	return ws.server.ListenAndServe(ws.config.Addr)
}
