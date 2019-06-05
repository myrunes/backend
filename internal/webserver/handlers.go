package webserver

import (
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/myrunes/internal/database"
	"github.com/zekroTJA/myrunes/internal/objects"
)

func (ws *WebServer) handlerFiles(ctx *routing.Context) error {
	path := string(ctx.Path())

	if strings.HasPrefix(path, "/api") {
		ctx.Next()
		return nil
	}

	if strings.HasPrefix(path, "/assets") ||
		strings.HasPrefix(path, "/css") ||
		strings.HasPrefix(path, "/js") ||
		strings.HasPrefix(path, "/favicon.ico") {

		fileHandlerStatic.NewRequestHandler()(ctx.RequestCtx)
		ctx.Abort()
		return nil
	}

	ctx.SendFile("./web/dist/index.html")
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

	if err = newUser.Validate(false); err != nil {
		return jsonResponse(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = ws.db.CreateUser(newUser); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	newUser.PassHash = nil

	ws.auth.CreateSession(ctx, newUser.UID, data.Remember)

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

func (ws *WebServer) handlerPostMe(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	reqUser := new(userRequest)
	var err error

	if err = parseJSONBody(ctx, reqUser); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if !ws.auth.CheckHash(user.PassHash, []byte(reqUser.CurrentPassword)) {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	newUser := &objects.User{
		UID:         user.UID,
		Username:    reqUser.Username,
		DisplayName: reqUser.DisplayName,
	}

	if err = newUser.Validate(true); err != nil {
		return jsonResponse(ctx, err, fasthttp.StatusBadRequest)
	}

	if reqUser.NewPassword != "" {
		newUser.PassHash, err = ws.auth.CreateHash([]byte(reqUser.NewPassword))
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	if _, err = ws.db.EditUser(newUser, false); err != nil {
		if err == database.ErrUsernameTaken {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerDeleteMe(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	reqUser := new(userRequest)
	var err error

	if err = parseJSONBody(ctx, reqUser); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if !ws.auth.CheckHash(user.PassHash, []byte(reqUser.CurrentPassword)) {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	if err = ws.db.DeleteUser(user.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerCreatePage(ctx *routing.Context) error {
	var err error
	page := objects.NewEmptyPage()

	if err = parseJSONBody(ctx, page); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = page.Validate(); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	user := ctx.Get("user").(*objects.User)
	page.FinalizeCreate(user.UID)

	if err = ws.db.CreatePage(page); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, page, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerGetPages(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	pages, err := ws.db.GetPages(user.UID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, &listResponse{N: len(pages), Data: pages}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetPage(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	page, err := ws.db.GetPage(uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if page == nil || page.Owner != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	return jsonResponse(ctx, page, fasthttp.StatusOK)
}

func (ws *WebServer) handlerEditPage(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	page, err := ws.db.GetPage(uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if page == nil || page.Owner != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if err = parseJSONBody(ctx, page); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	newPage, err := ws.db.EditPage(page)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, newPage, fasthttp.StatusOK)
}

func (ws *WebServer) handlerDeletePage(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	page, err := ws.db.GetPage(uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if page == nil || page.Owner != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if err = ws.db.DeletePage(page.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetChamps(ctx *routing.Context) error {
	return jsonResponse(ctx, &listResponse{N: len(objects.Champs), Data: objects.Champs}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetRunes(ctx *routing.Context) error {
	data := map[string]interface{}{
		"primary":   objects.RunesPrimary,
		"secondary": objects.RunesSecondary,
		"perks":     objects.PerksPool,
	}
	return jsonResponse(ctx, data, fasthttp.StatusOK)
}

func (ws *WebServer) handlerCheckUsername(ctx *routing.Context) error {
	uname := ctx.Param("uname")

	user, err := ws.db.GetUser(snowflake.ID(-1), strings.ToLower(uname))
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	status := fasthttp.StatusOK
	if user == nil {
		status = fasthttp.StatusNotFound
	}

	return jsonResponse(ctx, nil, status)
}
