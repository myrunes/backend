package webserver

import (
	"errors"
	"time"

	"github.com/zekroTJA/timedmap"

	"github.com/myrunes/backend/internal/assets"
	"github.com/myrunes/backend/internal/caching"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/mailserver"
	"github.com/myrunes/backend/internal/ratelimit"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// Error Objects
var (
	errNotFound                 = errors.New("not found")
	errInvalidArguments         = errors.New("invalid arguments")
	errUNameInUse               = errors.New("user name already in use")
	errNoAccess                 = errors.New("access denied")
	errMissingReCaptchaResponse = errors.New("missing recaptcha challenge response")
	errEmailAlreadyTaken        = errors.New("e-mail address is already taken by another account")
)

// Config wraps properties for the
// HTTP REST API server.
type Config struct {
	Addr       string           `json:"addr"`
	PathPrefix string           `json:"pathprefix"`
	TLS        *TLSConfig       `json:"tls"`
	ReCaptcha  *ReCaptchaConfig `json:"recaptcha"`
	PublicAddr string           `json:"publicaddress"`
	EnableCors bool             `json:"enablecors"`
	JWTKey     string           `json:"jwtkey"`
}

// TLSConfig wraps properties for
// TLS encryption.
type TLSConfig struct {
	Enabled bool   `json:"enabled"`
	Cert    string `json:"certfile"`
	Key     string `json:"keyfile"`
}

// ReCaptchaConfig wraps key and secret
// for ReCAPTCHA v2.
type ReCaptchaConfig struct {
	SiteKey   string `json:"sitekey"`
	SecretKey string `json:"secretkey"`
}

// WebServer provices a HTTP REST
// API router.
type WebServer struct {
	server *fasthttp.Server
	router *routing.Router

	db    database.Middleware
	cache caching.CacheMiddleware
	ms    *mailserver.MailServer
	auth  *Authorization
	rlm   *ratelimit.RateLimitManager

	avatarAssetsHandler *assets.AvatarHandler

	mailConfirmation *timedmap.TimedMap
	pwReset          *timedmap.TimedMap

	config *Config
}

// NewWebServer initializes a WebServer instance using
// the specified database driver, cache driver, mail
// server instance and configuration instance.
func NewWebServer(db database.Middleware, cache caching.CacheMiddleware,
	ms *mailserver.MailServer, avatarAssetsHandler *assets.AvatarHandler,
	config *Config) (ws *WebServer, err error) {

	ws = new(WebServer)

	ws.config = config
	ws.db = db
	ws.cache = cache
	ws.ms = ms
	ws.rlm = ratelimit.New()
	ws.router = routing.New()
	ws.server = &fasthttp.Server{
		Handler: ws.router.HandleRequest,
	}

	ws.avatarAssetsHandler = avatarAssetsHandler

	if ws.auth, err = NewAuthorization([]byte(config.JWTKey), db, cache, ws.rlm); err != nil {
		return
	}

	ws.mailConfirmation = timedmap.New(1 * time.Hour)
	ws.pwReset = timedmap.New(1 * time.Minute)

	ws.registerHandlers()

	return
}

// registerHandlers creates all rate limiter buckets and
// registers all routes and request handlers.
func (ws *WebServer) registerHandlers() {
	rlGlobal := ws.rlm.GetHandler(500*time.Millisecond, 50)
	rlUsersCreate := ws.rlm.GetHandler(15*time.Second, 1)
	rlPageCreate := ws.rlm.GetHandler(5*time.Second, 5)
	rlPostMail := ws.rlm.GetHandler(60*time.Second, 3)
	rlPwReset := ws.rlm.GetHandler(60*time.Second, 3)

	ws.router.Use(ws.addHeaders, rlGlobal)

	api := ws.router.Group(ws.config.PathPrefix)
	api.
		Post("/login", ws.handlerLogin)
	api.
		Post("/logout", ws.auth.CheckRequestAuth, ws.auth.Logout)

	api.Get("/version", ws.handlerGetVersion)
	api.Get("/recaptchainfo", ws.handlerGetReCaptchaInfo)

	assets := api.Group("/assets")
	assets.
		Get("/champions/avatars/<id>", ws.handlerGetAssetsChampionAvatars)

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

	pwReset := users.Group("/me/passwordreset")
	pwReset.
		Post("", rlPwReset, ws.handlerPostPwReset)
	pwReset.
		Post("/confirm", rlPwReset, ws.handlerPostPwResetConfirm)

	pages := api.Group("/pages", ws.addHeaders, rlGlobal, ws.auth.CheckRequestAuth)
	pages.
		Post("", rlPageCreate, ws.handlerCreatePage).
		Get(ws.handlerGetPages)
	pages.
		Get(`/<uid:\d+>`, ws.handlerGetPage).
		Post(ws.handlerEditPage).
		Delete(ws.handlerDeletePage)

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

// ListenAndServeBLocing starts the web servers
// listen and serving lifecycle which blocks
// the current goroutine.
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
