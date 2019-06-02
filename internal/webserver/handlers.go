package webserver

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/lol-runes/internal/objects"
)

func (ws *WebServer) handlerFiles(ctx *routing.Context) error {
	path := ctx.Path()
	if len(path) >= 4 && string(path[0:4]) == "/api" {
		ctx.Next()
		return nil
	}

	fileHandlerStatic.NewRequestHandler()(ctx.RequestCtx)
	ctx.Abort()
	return nil
}

func (ws *WebServer) handlerCreateUser(ctx *routing.Context) error {
	data := new(loginRequest)
	if err := parseJSONBody(ctx, data); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if data.UserName == "" || data.Password == "" {
		return jsonError(ctx, errInvalidArguments, fasthttp.StatusBadRequest)
	}

	user, err := ws.db.GetUser(snowflake.ID(-1), strings.ToLower(data.UserName))
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if user != nil {
		return jsonError(ctx, errUNameInUse, fasthttp.StatusConflict)
	}

	newUser, err := objects.NewUser(data.UserName, data.Password, ws.auth)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.CreateUser(newUser); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	newUser.PassHash = nil

	return jsonResponse(ctx, newUser, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerLogin(ctx *routing.Context) error {
	if !ws.auth.Login(ctx) {
		return nil
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetMe(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	user.PassHash = nil
	return jsonResponse(ctx, user, fasthttp.StatusOK)
}
