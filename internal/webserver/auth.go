package webserver

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"

	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
	"github.com/myrunes/backend/internal/static"
	"github.com/myrunes/backend/pkg/random"
	routing "github.com/qiangxue/fasthttp-routing"
	"golang.org/x/crypto/bcrypt"
)

const (
	attemptLimit          = 5 * time.Minute
	attemptBurst          = 5
	defCost               = 12
	apiTokenLength        = 64
	sessionKeyLength      = 128
	sessionExpireDefault  = 2 * time.Hour
	sessionExpireRemember = 30 * 24 * time.Hour

	jwtCookieName = "jwt_token"
)

var (
	errBadRequest   = errors.New("bad request")
	errUnauthorized = errors.New("unauthorized")
	errRateLimited  = errors.New("rate limited")

	setCookieHeader     = []byte("Set-Cookie")
	authorizationHeader = []byte("Authorization")

	jwtGenerationMethod = jwt.SigningMethodHS256
)

type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type Authorization struct {
	jwtKey []byte

	db  database.Middleware
	rlm *RateLimitManager
}

func NewAuthorization(jwtKey []byte, db database.Middleware, rlm *RateLimitManager) (auth *Authorization, err error) {
	auth = new(Authorization)
	auth.db = db
	auth.rlm = rlm

	if jwtKey == nil || len(jwtKey) == 0 {
		if auth.jwtKey, err = generateJWTKey(); err != nil {
			return
		}
	} else if len(jwtKey) < 32 {
		err = errors.New("JWT key must have at least 128 bit")
		return
	} else {
		auth.jwtKey = jwtKey
	}

	return
}

func (auth *Authorization) CreateHash(pass []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pass, defCost)
}

func (auth *Authorization) CheckHash(hash, pass []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, pass) == nil
}

func (auth *Authorization) CreateSessionKey() (string, error) {
	return random.Base64(sessionKeyLength)
}

func (auth *Authorization) Login(ctx *routing.Context) bool {
	login := new(loginRequest)
	if err := parseJSONBody(ctx, login); err != nil {
		return jsonError(ctx, errBadRequest, fasthttp.StatusBadRequest) != nil
	}

	limiter := auth.rlm.GetLimiter(fmt.Sprintf("loginAttempt#%s", getIPAddr(ctx)), attemptLimit, attemptBurst)

	if limiter.Tokens() <= 0 {
		return jsonError(ctx, errRateLimited, fasthttp.StatusTooManyRequests) != nil
	}

	user, err := auth.db.GetUser(snowflake.ID(-1), strings.ToLower(login.UserName))
	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError) != nil
	}
	if user == nil {
		limiter.Allow()
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized) != nil
	}

	if !auth.CheckHash(user.PassHash, []byte(login.Password)) {
		limiter.Allow()
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized) != nil
	}

	auth.CreateSession(ctx, user.UID, login.Remember)

	return true
}

func (auth *Authorization) CreateSession(ctx *routing.Context, uid snowflake.ID, remember bool) error {
	expires := time.Now()
	if remember {
		expires = expires.Add(sessionExpireRemember)
	} else {
		expires = expires.Add(sessionExpireDefault)
	}

	token, err := jwt.NewWithClaims(jwtGenerationMethod, jwt.StandardClaims{
		Subject:   uid.String(),
		ExpiresAt: expires.Unix(),
		IssuedAt:  time.Now().Unix(),
	}).SignedString(auth.jwtKey)

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	if _, err = auth.db.EditUser(&objects.User{UID: uid}, true); err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	secureCookie := ""
	if static.Release == "TRUE" {
		secureCookie = "; Secure"
	}

	cookie := fmt.Sprintf("%s=%s; Expires=%s; Path=/; HttpOnly%s",
		jwtCookieName, token, expires.Format(time.RFC1123), secureCookie)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	return nil
}

func (auth *Authorization) CheckRequestAuth(ctx *routing.Context) error {
	var user *objects.User
	var err error
	var apiToken string

	apiTokenB := ctx.Request.Header.PeekBytes(authorizationHeader)
	if apiTokenB != nil && len(apiTokenB) > 0 {
		apiToken := string(apiTokenB)
		if !strings.HasPrefix(strings.ToLower(apiToken), "basic ") {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		apiToken = apiToken[6:]
		user, err = auth.db.VerifyAPIToken(apiToken)
	} else {
		key := ctx.Request.Header.Cookie(jwtCookieName)
		if key == nil || len(key) == 0 {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		jwtToken, err := jwt.Parse(string(key), func(t *jwt.Token) (interface{}, error) {
			return auth.jwtKey, nil
		})
		if err != nil || !jwtToken.Valid {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		claims := jwt.StandardClaims{}
		claims.Subject, _ = claimsMap["sub"].(string)

		userID, _ := snowflake.ParseString(claims.Subject)
		user, err = auth.db.GetUser(userID, "")
	}

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if user == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	ctx.Set("user", user)
	ctx.Set("apitoken", apiToken)

	return nil
}

func (auth *Authorization) LogOut(ctx *routing.Context) error {
	key := ctx.Request.Header.Cookie("__session")
	if key == nil || len(key) == 0 {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	cookie := fmt.Sprintf("%s=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly", jwtCookieName)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func generateJWTKey() (key []byte, err error) {
	key = make([]byte, 32)
	_, err = rand.Read(key)
	return
}
