package webserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/myrunes/myrunes/pkg/comparison"
	"github.com/myrunes/myrunes/pkg/random"

	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/myrunes/myrunes/internal/database"
	"github.com/myrunes/myrunes/internal/objects"
	"github.com/myrunes/myrunes/internal/static"
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
	queryArgs := ctx.QueryArgs()

	sortBy := string(queryArgs.Peek("sortBy"))
	champion := string(queryArgs.Peek("champion"))
	short := string(queryArgs.Peek("short"))

	if champion == "" {
		champion = "general"
	}

	var sortFunc func(i, j *objects.Page) bool

	switch sortBy {
	case "created":
		sortFunc = func(i, j *objects.Page) bool {
			return i.Created.After(j.Created)
		}
	case "title":
		sortFunc = func(i, j *objects.Page) bool {
			return comparison.Alphabetically(i.Title, j.Title)
		}
	case "custom":
		if user.PageOrder != nil {
			pageOrder, ok := user.PageOrder[champion]
			if ok {
				sortFunc = func(i, j *objects.Page) bool {
					var pix, jix int
					for ix, uid := range pageOrder {
						if uid == i.UID {
							pix = ix
						} else if uid == j.UID {
							jix = ix
						}
					}
					return jix > pix
				}
			}
		}
	}

	pages, err := ws.db.GetPages(user.UID, champion, sortFunc)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if comparison.IsTrue(short) {
		m := make(map[string]int)
		for _, p := range pages {
			for _, c := range p.Champions {
				if _, ok := m[c]; !ok {
					m[c] = 1
				} else {
					m[c]++
				}
			}
		}

		return jsonResponse(ctx, &listResponse{N: len(m), Data: m}, fasthttp.StatusOK)
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
		"trees":     objects.RuneTrees,
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

func (ws *WebServer) handlerGetSessions(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	sessionKey := ctx.Get("sessionkey").(string)

	sessions, err := ws.db.GetSessions(user.UID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	var currentSession *objects.Session
	for _, s := range sessions {
		if sessionKey == s.Key {
			currentSession = s
		}
		s.Key = fmt.Sprintf("%s...%s", s.Key[:3], s.Key[len(s.Key)-3:])
	}

	return jsonResponse(ctx, sessionsResponse{
		listResponse: listResponse{
			N:    len(sessions),
			Data: sessions,
		},
		CurrentlyConnectedID: currentSession.SessionID.String(),
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerDeleteSession(ctx *routing.Context) error {
	_id := ctx.Param("uid")
	id, err := snowflake.ParseString(_id)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = ws.db.DeleteSession("", id); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostFavorite(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	var err error

	favReq := new(alterFavoriteRequest)
	if err = parseJSONBody(ctx, favReq); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if favReq.Favorites == nil {
		return jsonError(ctx, errBadRequest, fasthttp.StatusBadRequest)
	}

	champMap := make(map[string]interface{})
	for _, c := range objects.Champs {
		champMap[c] = nil
	}

	for i, f := range favReq.Favorites {
		favReq.Favorites[i] = strings.ToLower(f)
		if _, ok := champMap[f]; !ok {
			return jsonError(ctx, objects.ErrInvalidChamp, fasthttp.StatusBadRequest)
		}
	}

	user.Favorites = favReq.Favorites

	if _, err = ws.db.EditUser(user, false); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx,
		listResponse{N: len(user.Favorites), Data: user.Favorites},
		fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetFavorites(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	if user.Favorites == nil {
		user.Favorites = []string{}
	}

	return jsonResponse(ctx,
		listResponse{N: len(user.Favorites), Data: user.Favorites},
		fasthttp.StatusOK)
}

func (ws *WebServer) handlerCreateShare(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	var err error

	params := new(createShareRequest)
	if err := parseJSONBody(ctx, params); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	pageID, err := snowflake.ParseString(params.Page)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if page, err := ws.db.GetPage(pageID); err != nil {
		return jsonResponse(ctx, err, fasthttp.StatusInternalServerError)
	} else if page == nil || page.Owner != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	share, err := objects.NewSharePage(user.UID, pageID, params.MaxAccesses, params.Expires)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if err = ws.db.SetShare(share); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, share, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerPostShare(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	share, err := ws.db.GetShare("", uid, -1)
	if err != nil {
		return jsonResponse(ctx, err, fasthttp.StatusInternalServerError)
	}
	if share == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if share.OwnerID != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	params := new(createShareRequest)
	if err := parseJSONBody(ctx, params); err != nil {
		return jsonError(ctx, errBadRequest, fasthttp.StatusBadRequest)
	}

	if (params.Expires != time.Time{}) {
		share.Expires = params.Expires
	}

	if params.MaxAccesses != 0 {
		share.MaxAccesses = params.MaxAccesses
	}

	if err = ws.db.SetShare(share); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, share, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerGetShare(ctx *routing.Context) error {
	ident := ctx.Param("ident")
	byIdent := true
	var user *objects.User

	_user := ctx.Get("user")
	if _user != nil {
		user = _user.(*objects.User)
		byIdent = false
	}

	shareID, err := snowflake.ParseString(ident)
	if err != nil || user == nil {
		shareID = -1
	}

	share, err := ws.db.GetShare(ident, shareID, shareID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if share == nil || (user != nil && share.OwnerID != user.UID) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if byIdent && (share.MaxAccesses == 0 || (share.Expires != time.Time{} && share.Expires.Before(time.Now()))) {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	page, err := ws.db.GetPage(share.PageID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if page == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	owner, err := ws.db.GetUser(page.Owner, "")
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	reqAddr := getIPAddr(ctx)
	validReqAddr := !strings.HasPrefix(reqAddr, "192.168") &&
		!strings.HasPrefix(reqAddr, "10.23") &&
		!(static.Release == "TRUE" && reqAddr == "127.0.0.1") &&
		string(ctx.Request.Header.PeekBytes(headerUserAgent)) != static.DiscordUserAgentPingHeaderVal

	if byIdent && validReqAddr {
		var contains bool
		for _, ip := range share.AccessIPs {
			if ip == reqAddr {
				contains = true
			}
		}

		if !contains {
			share.AccessIPs = append(share.AccessIPs, reqAddr)

			if share.MaxAccesses > 0 {
				share.MaxAccesses--
			}

			share.LastAccess = time.Now()
			share.Accesses++

			if err = ws.db.SetShare(share); err != nil {
				return jsonError(ctx, err, fasthttp.StatusInternalServerError)
			}
		}
	}

	share.AccessIPs = nil

	return jsonResponse(ctx, &shareResponse{
		Page:  page,
		Share: share,
		User:  owner,
	}, fasthttp.StatusAccepted)
}

func (ws *WebServer) handlerDeleteShare(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	share, err := ws.db.GetShare("", uid, -1)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if share == nil || share.OwnerID != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if err = ws.db.DeleteShare("", uid, -1); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetVersion(ctx *routing.Context) error {
	return jsonResponse(ctx, map[string]string{
		"version": static.AppVersion,
		"release": static.Release,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostAPIToken(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	var err error
	token := new(objects.APIToken)

	if token.Token, err = random.GetRandBase64Str(apiTokenLength); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	token.UserID = user.UID
	token.Created = time.Now()

	if err = ws.db.SetAPIToken(token); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, token, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetAPIToken(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	token, err := ws.db.GetAPIToken(user.UID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if token == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	return jsonResponse(ctx, token, fasthttp.StatusOK)
}

func (ws *WebServer) handlerDeleteAPIToken(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	if err := ws.db.ResetAPIToken(user.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostPageOrder(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	queryArgs := ctx.QueryArgs()
	champion := string(queryArgs.Peek("champion"))

	if champion == "" {
		champion = "general"
	}

	pageOrder := new(pageOrderRequest)
	if err := parseJSONBody(ctx, pageOrder); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if user.PageOrder == nil {
		user.PageOrder = make(map[string][]snowflake.ID)
	}

	user.PageOrder[champion] = pageOrder.PageOrder
	if _, err := ws.db.EditUser(user, false); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}
