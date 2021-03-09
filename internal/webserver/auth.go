package webserver

import (
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
	// time until a new rate limiter ticket is
	// generated for login tries
	attemptLimit = 5 * time.Minute
	// ammount of tickets which can be stashed
	// for login attempts
	attemptBurst = 5
	// character length of generated API tokens
	apiTokenLength = 64
	// The byte length of the signing key of
	// accessTokens
	signingKeyLength = 128
	// The character length of refreshTokens
	refreshTokenLength = 64
	// defaul time until a default login
	// session expires
	sessionExpireDefault = 2 * time.Hour
	// default time until a "remembered"
	// login session expires
	sessionExpireRemember = 30 * 24 * time.Hour
	// time until an access token must be
	// reacquired
	accessTokenLifetime = 1 * time.Hour
	// cookie key name of the refreshToken
	refreshTokenCookieName = "refreshToken"
)

var (
	errBadRequest    = errors.New("bad request")
	errUnauthorized  = errors.New("unauthorized")
	errInvalidAccess = errors.New("invalid access key")
	errRateLimited   = errors.New("rate limited")

	setCookieHeader     = []byte("Set-Cookie")
	authorizationHeader = []byte("Authorization")

	jwtGenerationMethod = jwt.SigningMethodHS256

	argon2Params = getArgon2Params()
)

// loginRequests describes the request
// model of the login endpoint
type loginRequest struct {
	reCaptchaResponse

	UserName string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// Authorization provides functionalities
// for HTTP session authorization and
// session lifecycle maintainance.
type Authorization struct {
	signingKey []byte

	db    database.Middleware
	cache caching.CacheMiddleware
	rlm   *ratelimit.RateLimitManager
}

// NewAuthorization initializes a new
// Authorization instance using the passed
// jwtKey, which will be used to sign JWTs,
// the database driver, cache driver and
// rate limit manager.
// If the passed jwtKey is nil or empty,
// a random key will be generated on
// initialization.
func NewAuthorization(signingKey []byte, db database.Middleware, cache caching.CacheMiddleware, rlm *ratelimit.RateLimitManager) (auth *Authorization, err error) {
	auth = new(Authorization)
	auth.db = db
	auth.cache = cache
	auth.rlm = rlm

	if signingKey == nil || len(signingKey) == 0 {
		if auth.signingKey, err = random.ByteArray(signingKeyLength); err != nil {
			return
		}
	} else if len(signingKey) < 32 {
		err = errors.New("JWT key must have at least 128 bit")
		return
	} else {
		auth.signingKey = signingKey
	}

	return
}

// CreateHash creates a hash string from the passed
// pass string containing information about the used
// algorithm and parameters used to generate the hash
// together with the actual hash data.
//
// This implementation uses Argon2id hash generation.
func (auth *Authorization) CreateHash(pass string) (string, error) {
	return argon2id.CreateHash(pass, argon2Params)
}

// CheckHash tries to compare the passed hash string
// with the passed pass string by using the method and
// parameters specified in the hash string.
//
// This imlementation supports both the old hash
// algorithm used in myrunes before batch 1.7.x
// (bcrypt) and the current implementation argon2id.
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

// Login provides a handler accepting login credentials
// as JSON POST body. This is used to authenticate a user
// and create a login session on successful authentication.
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

	if token, err := auth.CreateAndSetRefreshToken(ctx, user.UID, login.Remember); err != nil {
		auth.cache.SetUserByToken(token, user)
	}

	return true
}

// CreateSession creates a login session for the specified
// user. This generates a JWT which is signed with the internal
// jwtKey and then stored as cookie on response.
func (auth *Authorization) CreateAndSetRefreshToken(ctx *routing.Context, uid snowflake.ID, remember bool) (token string, err error) {
	expires := time.Now()
	if remember {
		expires = expires.Add(sessionExpireRemember)
	} else {
		expires = expires.Add(sessionExpireDefault)
	}

	if token, err = random.Base64(refreshTokenLength); err != nil {
		err = jsonError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	err = auth.db.SetRefreshToken(&objects.RefreshToken{
		Token:    token,
		UserID:   uid,
		Deadline: expires,
	})
	if err != nil {
		err = jsonError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	user, err := auth.cache.GetUserByID(uid)
	if err != nil {
		return
	}

	user.Update(nil, true)
	if err = user.Validate(true); err != nil {
		err = jsonError(ctx, err, fasthttp.StatusBadRequest)
		return
	}
	if err = auth.db.EditUser(user); err != nil {
		err = jsonError(ctx, err, fasthttp.StatusInternalServerError)
		return
	}

	cookieSecurity := ""
	if static.Release == "TRUE" {
		cookieSecurity = "; Secure; SameSite=Strict"
	}

	cookie := fmt.Sprintf("%s=%s; Expires=%s; Path=/; HttpOnly%s",
		refreshTokenCookieName, token, expires.Format(time.RFC1123), cookieSecurity)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	return
}

func (auth *Authorization) ObtainAccessToken(ctx *routing.Context) (string, error) {
	key := ctx.Request.Header.Cookie(refreshTokenCookieName)
	if key == nil || len(key) == 0 {
		return "", jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	refreshToken := string(key)
	token, err := auth.db.GetRefreshToken(refreshToken)
	if err != nil {
		return "", jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if token == nil {
		return "", jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	now := time.Now()
	if now.After(token.Deadline) {
		auth.db.RemoveRefreshToken(token.Token)
		return "", jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	accessToken, err := jwt.NewWithClaims(jwtGenerationMethod, jwt.StandardClaims{
		Subject:   token.UserID.String(),
		ExpiresAt: now.Add(accessTokenLifetime).Unix(),
		IssuedAt:  now.Unix(),
	}).SignedString(auth.signingKey)
	if err != nil {
		return "", jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}

	return accessToken, nil
}

// CheckRequestAuth provides a handler which
// cancels the current handler stack if no valid
// session authentication or API token could be
// identified in the request.
func (auth *Authorization) CheckRequestAuth(ctx *routing.Context) error {
	var user *objects.User
	var err error
	var authValue string

	authValueB := ctx.Request.Header.PeekBytes(authorizationHeader)
	if authValueB != nil && len(authValueB) > 0 {
		authValue = string(authValueB)
	}
	if strings.HasPrefix(strings.ToLower(authValue), "basic") {
		authValue = authValue[6:]
		var ok bool
		if user, ok = auth.cache.GetUserByToken(authValue); !ok {
			if user, err = auth.db.VerifyAPIToken(authValue); err == nil {
				auth.cache.SetUserByToken(authValue, user)
			}
		}
	} else if strings.HasPrefix(strings.ToLower(authValue), "accesstoken ") {
		authValue = authValue[12:]

		jwtToken, err := jwt.Parse(authValue, func(t *jwt.Token) (interface{}, error) {
			return auth.signingKey, nil
		})
		if err != nil || !jwtToken.Valid {
			return jsonError(ctx, errInvalidAccess, fasthttp.StatusUnauthorized)
		}

		claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return jsonError(ctx, errInvalidAccess, fasthttp.StatusUnauthorized)
		}

		claims := jwt.StandardClaims{}
		claims.Subject, _ = claimsMap["sub"].(string)

		userID, _ := snowflake.ParseString(claims.Subject)
		user, err = auth.cache.GetUserByID(userID)
	}

	if err != nil {
		return jsonError(ctx, err, fasthttp.StatusInternalServerError)
	}
	if user == nil {
		return jsonError(ctx, errInvalidAccess, fasthttp.StatusUnauthorized)
	}

	ctx.Set("user", user)
	ctx.Set("apitoken", authValue)
	// ctx.Set("jwt", jwtTokenStr)

	return nil
}

// Logout provides a handler which removes the
// session JWT cookie by setting an invalid,
// expired session cookie.
func (auth *Authorization) Logout(ctx *routing.Context) error {
	key := ctx.Request.Header.Cookie(refreshTokenCookieName)
	if key == nil || len(key) == 0 {
		return jsonError(ctx, errUnauthorized, fasthttp.StatusUnauthorized)
	}

	cookie := fmt.Sprintf("%s=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Path=/; HttpOnly", refreshTokenCookieName)
	ctx.Response.Header.AddBytesK(setCookieHeader, cookie)

	// if jwt, ok := ctx.Get("jwt").(string); ok {
	// 	auth.cache.SetUserByToken(jwt, nil)
	// }

	return jsonResponse(ctx, nil, fasthttp.StatusOK)
}

// getArgon2Params returns an instance of default
// parameters which are used for generating
// Argon2id password hashes.
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
