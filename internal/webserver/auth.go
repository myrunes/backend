package webserver

import (
	"crypto/rand"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"

	"github.com/myrunes/backend/internal/caching"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
	"github.com/myrunes/backend/internal/ratelimit"
	"github.com/myrunes/backend/internal/shared"
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

	argon2Params = getArgon2Params()
)

type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type Authorization struct {
	jwtKey []byte

	db    database.Middleware
	cache caching.Middleware
	rlm   *ratelimit.RateLimitManager
}

func NewAuthorization(jwtKey []byte, db database.Middleware, cache caching.Middleware, rlm *ratelimit.RateLimitManager) (auth *Authorization, err error) {
	auth = new(Authorization)
	auth.db = db
	auth.cache = cache
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

func (auth *Authorization) CreateHash(pass string) (string, error) {
	return argon2id.CreateHash(pass, argon2Params)
}

func (auth *Authorization) CheckHash(hash, pass string) bool {
	if strings.HasPrefix(hash, "$2a") {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
	}

	if strings.HasPrefix(hash, "$argon2id") {
		ok, err := argon2id.ComparePasswordAndHash(pass, hash)
		return ok && err == nil
	}

	return false
}

func (auth *Authorization) CreateSessionKey() (string, error) {
	return random.Base64(sessionKeyLength)
}

func (auth *Authorization) Login(ctx *routing.Context) bool {
	login := new(loginRequest)
	if err := parseJSONBody(ctx, login); err != nil {
		return jsonError(ctx, errBadRequest, fasthttp.StatusBadRequest) != nil
	}

	limiter := auth.rlm.GetLimiter(fmt.Sprintf("loginAttempt#%s", shared.GetIPAddr(ctx)), attemptLimit, attemptBurst)

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

	// Querrying user in cache to set cache entry
	auth.cache.GetUserByID(user.UID)

	if !auth.CheckHash(string(user.PassHash), login.Password) {
		limiter.Allow()
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized) != nil
	}

	if token, err := auth.CreateSession(ctx, user.UID, login.Remember); err != nil {
		auth.cache.SetUserByToken(token, user)
	}

	return true
}

func (auth *Authorization) CreateSession(ctx *routing.Context, uid snowflake.ID, remember bool) (string, error) {
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
		return "", jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	user, err := auth.cache.GetUserByID(uid)
	if err != nil {
		return "", err
	}

	user.Update(nil, true)
	if err = user.Validate(true); err != nil {
		return "", jsonError(ctx, err, fasthttp.StatusBadRequest)
	}

	if err = auth.db.EditUser(user); err != nil {
		return "", jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	secureCookie := ""
	if static.Release == "TRUE" {
		secureCookie = "; Secure; SameSite=Strict"
	}

	cookie := fmt.Sprintf("%s=%s; Expires=%s; Path=/; HttpOnly%s",
		jwtCookieName, token, expires.Format(time.RFC1123), secureCookie)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	return token, nil
}

func (auth *Authorization) CheckRequestAuth(ctx *routing.Context) error {
	var user *objects.User
	var err error
	var jwtTokenStr string
	var apiToken string

	apiTokenB := ctx.Request.Header.PeekBytes(authorizationHeader)
	if apiTokenB != nil && len(apiTokenB) > 0 {
		apiToken := string(apiTokenB)
		if !strings.HasPrefix(strings.ToLower(apiToken), "basic ") {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		apiToken = apiToken[6:]
		var ok bool
		if user, ok = auth.cache.GetUserByToken(apiToken); !ok {
			if user, err = auth.db.VerifyAPIToken(apiToken); err == nil {
				auth.cache.SetUserByToken(apiToken, user)
			}
		}
	} else {
		key := ctx.Request.Header.Cookie(jwtCookieName)
		if key == nil || len(key) == 0 {
			return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
		}

		jwtTokenStr = string(key)

		var ok bool
		if user, ok = auth.cache.GetUserByToken(jwtTokenStr); !ok {
			jwtToken, err := jwt.Parse(jwtTokenStr, func(t *jwt.Token) (interface{}, error) {
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
			user, err = auth.cache.GetUserByID(userID)

			auth.cache.SetUserByToken(jwtTokenStr, user)
		}

	}

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if user == nil {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	ctx.Set("user", user)
	ctx.Set("apitoken", apiToken)
	ctx.Set("jwt", jwtTokenStr)

	return nil
}

func (auth *Authorization) LogOut(ctx *routing.Context) error {
	key := ctx.Request.Header.Cookie(jwtCookieName)
	if key == nil || len(key) == 0 {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	cookie := fmt.Sprintf("%s=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly", jwtCookieName)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	if jwt, ok := ctx.Get("jwt").(string); ok {
		auth.cache.SetUserByToken(jwt, nil)
	}

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

func generateJWTKey() (key []byte, err error) {
	key = make([]byte, 32)
	_, err = rand.Read(key)
	return
}

func getArgon2Params() *argon2id.Params {
	cpus := runtime.NumCPU()

	return &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: uint8(cpus),
		SaltLength:  16,
		KeyLength:   32,
	}
}
