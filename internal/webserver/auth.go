package webserver

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/valyala/fasthttp"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/zekroTJA/lol-runes/internal/database"
	"github.com/zekroTJA/lol-runes/internal/objects"
	"golang.org/x/crypto/bcrypt"
)

const (
	defCost               = 12
	sessionKeyLength      = 128
	sessionExpireDefault  = 2 * time.Hour
	sessionExpireRemember = 30 * 24 * time.Hour
)

var (
	errBadRequest   = errors.New("bad request")
	errUnauthorized = errors.New("unauthorized")

	setCookieHeader = []byte("Set-Cookie")
)

type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type Authorization struct {
	db database.Middleware
}

func NewAuthorization(db database.Middleware) (auth *Authorization) {
	auth = new(Authorization)
	auth.db = db
	return
}

func (auth *Authorization) CreateHash(pass []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pass, defCost)
}

func (auth *Authorization) CheckHash(hash, pass []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, pass) == nil
}

func (auth *Authorization) CreateSessionKey() (string, error) {
	key := make([]byte, sessionKeyLength)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func (auth *Authorization) Login(ctx *routing.Context) bool {
	login := new(loginRequest)
	if err := parseJSONBody(ctx, login); err != nil {
		return jsonError(ctx, errBadRequest, fasthttp.StatusBadRequest) != nil
	}

	user, err := auth.db.GetUser(snowflake.ID(-1), strings.ToLower(login.UserName))
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError) != nil
	}
	if user == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized) != nil
	}

	if !auth.CheckHash(user.PassHash, []byte(login.Password)) {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized) != nil
	}

	sessionKey, err := auth.CreateSessionKey()
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError) != nil
	}

	expires := time.Now()
	if login.Remember {
		expires = expires.Add(sessionExpireRemember)
	} else {
		expires = expires.Add(sessionExpireDefault)
	}

	if err = auth.db.CreateSession(sessionKey, user.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError) != nil
	}

	if _, err = auth.db.EditUser(&objects.User{UID: user.UID}, true); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError) != nil
	}

	cookie := fmt.Sprintf("__session=%s; Expires=%s; Path=/; HttpOnly",
		sessionKey, expires.Format(time.RFC3339))
	ctx.Response.Header.SetBytesK(setCookieHeader, cookie)

	return true
}

func (auth *Authorization) CheckRequestAuth(ctx *routing.Context) error {
	key := ctx.Request.Header.Cookie("__session")
	if key == nil || len(key) == 0 {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	user, err := auth.db.GetSession(string(key))
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if user == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	ctx.Set("user", user)

	return nil
}
