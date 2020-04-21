package webserver

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/myrunes/backend/internal/ddragon"
	"github.com/myrunes/backend/internal/shared"

	"github.com/myrunes/backend/pkg/comparison"
	"github.com/myrunes/backend/pkg/random"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
	"github.com/myrunes/backend/internal/static"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func (ws *WebServer) handlerCreateUser(ctx *routing.Context) error {
	data := new(loginRequest)
	if err := parseJSONBody(ctx, data); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if data.UserName == "" || data.Password == "" || len(data.Password) < 8 {
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

	ws.auth.CreateSession(ctx, newUser.UID, data.Remember)

	outUser := *newUser
	outUser.PassHash = nil

	return jsonResponse(ctx, outUser, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerLogin(ctx *routing.Context) error {
	if !ws.auth.Login(ctx) {
		return nil
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetMe(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	userOut := *user
	userOut.PassHash = nil
	return jsonResponse(ctx, userOut, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostMe(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	reqUser := new(userRequest)
	var err error

	if err = parseJSONBody(ctx, reqUser); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if !ws.auth.CheckHash(string(user.PassHash), reqUser.CurrentPassword) {
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
		if len(reqUser.NewPassword) < 8 {
			return jsonError(ctx, fmt.Errorf("invalid new password"), fasthttp.StatusBadRequest)
		}
		var passStr string
		passStr, err = ws.auth.CreateHash(reqUser.NewPassword)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
		newUser.PassHash = []byte(passStr)
	}

	if _, err = ws.db.EditUser(newUser, false); err != nil {
		if err == database.ErrUsernameTaken {
			return jsonError(ctx, err, fasthttp.StatusBadRequest)
		}
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	ws.cache.SetUserByID(newUser.UID, newUser)
	if jwtToken, ok := ctx.Get("jwt").(string); ok {
		ws.cache.SetUserByJWT(jwtToken, newUser)
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

	if !ws.auth.CheckHash(string(user.PassHash), reqUser.CurrentPassword) {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	if err = ws.db.DeleteUser(user.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	ws.cache.SetUserByID(user.UID, nil)
	if jwtToken, ok := ctx.Get("jwt").(string); ok {
		ws.cache.SetUserByJWT(jwtToken, nil)
	}

	return ws.auth.LogOut(ctx)
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
	ws.cache.SetPageByID(page.UID, page)

	return jsonResponse(ctx, page, fasthttp.StatusCreated)
}

func (ws *WebServer) handlerGetPages(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	queryArgs := ctx.QueryArgs()

	sortBy := string(queryArgs.Peek("sortBy"))
	filter := string(queryArgs.Peek("filter"))
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

	pages, err := ws.db.GetPages(user.UID, champion, filter, sortFunc)
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

	page, err := ws.cache.GetPageByID(uid)
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

	page, err := ws.cache.GetPageByID(uid)
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
	ws.cache.SetPageByID(newPage.UID, newPage)

	return jsonResponse(ctx, newPage, fasthttp.StatusOK)
}

func (ws *WebServer) handlerDeletePage(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	_uid := ctx.Param("uid")
	uid, err := snowflake.ParseString(_uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	page, err := ws.cache.GetPageByID(uid)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if page == nil || page.Owner != user.UID {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	if err = ws.db.DeletePage(page.UID); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	ws.cache.SetPageByID(page.UID, nil)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetChamps(ctx *routing.Context) error {
	return jsonCachableResponse(ctx, &listResponse{N: len(ddragon.DDragonInstance.Champions), Data: ddragon.DDragonInstance.Champions}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerGetRunes(ctx *routing.Context) error {
	data := map[string]interface{}{
		"trees": ddragon.DDragonInstance.Runes,
		"perks": objects.PerksPool,
	}
	return jsonCachableResponse(ctx, data, fasthttp.StatusOK)
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

// TODO: DEPRECATED -- REMOVE
func (ws *WebServer) handlerGetSessions(ctx *routing.Context) error {
	return jsonError(ctx, errors.New("deprecated"), fasthttp.StatusGone)
}

// TODO: DEPRECATED -- REMOVE
func (ws *WebServer) handlerDeleteSession(ctx *routing.Context) error {
	return jsonError(ctx, errors.New("deprecated"), fasthttp.StatusGone)
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
	for _, c := range ddragon.DDragonInstance.Champions {
		champMap[c.UID] = nil
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

	if page, err := ws.cache.GetPageByID(pageID); err != nil {
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

	page, err := ws.cache.GetPageByID(share.PageID)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if page == nil {
		return jsonError(ctx, errNotFound, fasthttp.StatusNotFound)
	}

	owner, err := ws.cache.GetUserByID(page.Owner)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	reqAddr := shared.GetIPAddr(ctx)
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
	return jsonCachableResponse(ctx, map[string]string{
		"version":    static.AppVersion,
		"apiversion": static.APIVersion,
		"release":    static.Release,
	}, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostAPIToken(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)
	var err error
	token := new(objects.APIToken)

	if token.Token, err = random.Base64(apiTokenLength); err != nil {
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

func (ws *WebServer) handlerPostMail(ctx *routing.Context) error {
	user := ctx.Get("user").(*objects.User)

	mail := new(setMailRequest)
	if err := parseJSONBody(ctx, mail); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if !ws.auth.CheckHash(string(user.PassHash), mail.CurrentPassword) {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	if mail.Reset {
		_, err := ws.db.EditUser(&objects.User{
			UID:         user.UID,
			MailAddress: "__RESET__",
		}, false)
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}

		return jsonResponse(ctx, nil, fasthttp.StatusOK)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token, err := random.String(16, charset)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	mailText := fmt.Sprintf(
		"Please open the following link to confirm your E-Mail address:\n"+
			"%s/mailConfirmation?token=%s", ws.config.PublicAddr, token)

	err = ws.ms.SendMailFromDef(mail.MailAddress, "E-Mail confirmation | myrunes", mailText, "text/plain")
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	ws.mailConfirmation.Set(token, &mailConfirmationData{
		MailAddress: mail.MailAddress,
		UserID:      user.UID,
	}, 12*time.Hour)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostConfirmMail(ctx *routing.Context) error {
	token := new(confirmMail)
	if err := parseJSONBody(ctx, token); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if !ws.mailConfirmation.Contains(token.Token) {
		return jsonError(ctx, fmt.Errorf("invalid token"), fasthttp.StatusBadRequest)
	}

	data, ok := ws.mailConfirmation.GetValue(token.Token).(*mailConfirmationData)
	if !ok {
		return jsonError(ctx, fmt.Errorf("wrong data struct in timedmap"), fasthttp.StatusInternalServerError)
	}

	ws.mailConfirmation.Remove(token.Token)

	_, err := ws.db.EditUser(&objects.User{
		UID:         data.UserID,
		MailAddress: data.MailAddress,
	}, false)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, err, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostPwReset(ctx *routing.Context) error {
	reset := new(passwordReset)
	if err := parseJSONBody(ctx, reset); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	user, err := ws.db.GetUser(-1, reset.MailAddress)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if user == nil || user.MailAddress == "" {
		return jsonResponse(ctx, nil, fasthttp.StatusOK)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token, err := random.String(24, charset)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	mailText := fmt.Sprintf("Please follow the link below to reset your accounts password:\n"+
		"%s/passwordReset?token=%s", ws.config.PublicAddr, token)
	err = ws.ms.SendMailFromDef(user.MailAddress, "Password reset | myrunes", mailText, "text/plain")
	if err == nil {
		ws.pwReset.Set(token, user.UID, 10*time.Minute)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func (ws *WebServer) handlerPostPwResetConfirm(ctx *routing.Context) error {
	data := new(confirmPasswordReset)
	if err := parseJSONBody(ctx, data); err != nil {
		return jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if data.NewPassword == "" || len(data.NewPassword) < 8 {
		return jsonError(ctx, fmt.Errorf("invalid password length"), fasthttp.StatusBadRequest)
	}

	if !ws.pwReset.Contains(data.Token) {
		return jsonError(ctx, fmt.Errorf("invalid token"), fasthttp.StatusBadRequest)
	}

	uID, ok := ws.pwReset.GetValue(data.Token).(snowflake.ID)
	if !ok {
		return jsonError(ctx, fmt.Errorf("wrong data struct in timedmap"), fasthttp.StatusInternalServerError)
	}

	user, err := ws.db.GetUser(uID, "")
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if user == nil {
		return jsonError(ctx, fmt.Errorf("unknown user"), fasthttp.StatusBadRequest)
	}

	errCheckFailed := fmt.Errorf("security check failed")
	if len(data.PageNames) < 3 || data.PageNames[0] == "" || data.PageNames[1] == "" || data.PageNames[2] == "" {
		return jsonError(ctx, errCheckFailed, fasthttp.StatusBadRequest)
	}

	pages, err := ws.db.GetPages(uID, "", "", nil)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	checkMap := make(map[string]interface{})
	for _, guess := range data.PageNames {
		if _, ok := checkMap[guess]; ok {
			return jsonError(ctx, errCheckFailed, fasthttp.StatusBadRequest)
		}
		checkMap[guess] = nil
	}

	var guessed int

	for _, page := range pages {
		for i, guess := range data.PageNames {
			if checkPageName(page.Title, guess, 0.2) {
				guessed++
				data.PageNames[i] = ""
			}
		}
	}

	if guessed < 3 {
		return jsonError(ctx, errCheckFailed, fasthttp.StatusBadRequest)
	}

	newUser := &objects.User{
		UID: user.UID,
	}

	ws.pwReset.Remove(data.Token)

	var passStr string
	passStr, err = ws.auth.CreateHash(data.NewPassword)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	newUser.PassHash = []byte(passStr)

	_, err = ws.db.EditUser(newUser, false)
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}
