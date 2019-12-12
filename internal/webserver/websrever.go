package webserver

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/timedmap"

	"github.com/myrunes/myrunes/internal/database"
	"github.com/myrunes/myrunes/internal/mailserver"

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
		Compress:   true,
		// PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
		// 	return ctx.Path()[7:]
		// },
	}
)

type Config struct {
	Addr       string     `json:"addr"`
	TLS        *TLSConfig `json:"tls"`
	PublicAddr string     `json:"publicaddress"`
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
	ms   *mailserver.MailServer
	auth *Authorization
	rlm  *RateLimitManager

	mailConfirmation *timedmap.TimedMap

	config *Config
}

type mailConfirmationData struct {
	UserID      snowflake.ID
	MailAddress string
}

func NewWebServer(db database.Middleware, ms *mailserver.MailServer, config *Config, assets string) (ws *WebServer) {
	ws = new(WebServer)

	if assets != "" {
		fileHandlerStatic.Root = assets
	}

	ws.config = config
	ws.db = db
	ws.ms = ms
	ws.rlm = NewRateLimitManager()
	ws.auth = NewAuthorization(db, ws.rlm)
	ws.router = routing.New()
	ws.server = &fasthttp.Server{
		Handler: ws.router.HandleRequest,
	}

	ws.mailConfirmation = timedmap.New(1 * time.Hour)

	ws.registerHandlers()

	return
}

func (ws *WebServer) registerHandlers() {
	rlGlobal := ws.rlm.GetHandler(500*time.Millisecond, 50)
	rlUsersCreate := ws.rlm.GetHandler(15*time.Second, 1)
	rlPageCreate := ws.rlm.GetHandler(5*time.Second, 5)
	rlPostMail := ws.rlm.GetHandler(60*time.Second, 1)

	ws.router.Use(ws.handlerFiles, ws.addHeaders, rlGlobal)

	api := ws.router.Group("/api")
	api.
		Post("/login", ws.handlerLogin)
	api.
		Post("/logout", ws.auth.LogOut)

	api.Get("/version", ws.handlerGetVersion)

	resources := api.Group("/resources")
	resources.
		Get("/champions", ws.handlerGetChamps)
	resources.
		Get("/runes", ws.handlerGetRunes)

	users := api.Group("/users")
	users.
		Post("", rlUsersCreate, ws.handlerCreateUser)
	users.
		Post("/me", ws.auth.CheckRequestAuth, ws.handlerPostMe).
		Get(ws.auth.CheckRequestAuth, ws.handlerGetMe).
		Delete(ws.auth.CheckRequestAuth, ws.handlerDeleteMe)
	users.
		Get("/<uname>", ws.handlerCheckUsername)
	users.
		Post("/me/pageorder", ws.auth.CheckRequestAuth, ws.handlerPostPageOrder)

	email := users.Group("/me/mail")
	email.
		Post("", ws.auth.CheckRequestAuth, rlPostMail, ws.handlerPostMail)
	email.
		Post("/confirm", ws.handlerPostConfirmMail)

	pages := api.Group("/pages", ws.addHeaders, rlGlobal, ws.auth.CheckRequestAuth)
	pages.
		Post("", rlPageCreate, ws.handlerCreatePage).
		Get(ws.handlerGetPages)
	pages.
		Get(`/<uid:\d+>`, ws.handlerGetPage).
		Post(ws.handlerEditPage).
		Delete(ws.handlerDeletePage)

	sessions := api.Group("/sessions", ws.addHeaders, rlGlobal, ws.auth.CheckRequestAuth)
	sessions.
		Get("", ws.handlerGetSessions)
	sessions.
		Delete(`/<uid:\d+>`, ws.handlerDeleteSession)

	favorites := api.Group("/favorites", ws.addHeaders, rlGlobal, ws.auth.CheckRequestAuth)
	favorites.
		Get("", ws.handlerGetFavorites).
		Post(ws.handlerPostFavorite)

	shares := api.Group("/shares", ws.addHeaders, rlGlobal)
	shares.
		Post("", ws.auth.CheckRequestAuth, ws.handlerCreateShare)
	shares.
		Get(`/<ident:\d+>`, ws.auth.CheckRequestAuth, ws.handlerGetShare)
	shares.
		Get("/<ident:.+>", ws.handlerGetShare)
	shares.
		Post(`/<uid:\d+>`, ws.auth.CheckRequestAuth, ws.handlerPostShare).
		Delete(ws.auth.CheckRequestAuth, ws.handlerDeleteShare)

	apitoken := api.Group("/apitoken", ws.addHeaders, rlGlobal, ws.auth.CheckRequestAuth)
	apitoken.
		Get("", ws.handlerGetAPIToken).
		Post(ws.handlerPostAPIToken).
		Delete(ws.handlerDeleteAPIToken)

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
