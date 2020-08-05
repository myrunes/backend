package auth

import (
	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
)

// AuthMiddleware describes a module which provides
// functionality to hash passwords, check password
// hases, create HTTP sessions and authorize HTTP
// requests.
type AuthMiddleware interface {

	// CreateHash creates a secure password hash
	// for the given password string and returns the
	// hash as string containing the hashing algorithm,
	// the parameters used to create the hash and the
	// hash itself as base64 vlaue.
	CreateHash(pass string) (string, error)

	// CheckHash checks if the given hash matches a
	// given password. The result is returned as
	// boolean. If something fails during the hashing,
	// the reutrned  result will be 'false'.
	CheckHash(hash, pass string) bool

	// CreateSession creates a new web session for the
	// passed uid with a defined life span which is
	// expanded when remember is true. The session is
	// then set as request cookie to identify and
	// authorize following requests from this user.
	CreateSession(ctx *routing.Context, uid snowflake.ID, remember bool) (string, error)

	// Login collects login credentials from the
	// request payload. After successful authorization,
	// a session will be generated and set to the
	// response via CreateSession.
	// Otherwise, a 401 Untauthorized response will
	// be sent back.
	Login(ctx *routing.Context) bool

	// Logout removes the session identification
	// from the requested user so that following
	// requests can not be authorized anymore.
	Logout(ctx *routing.Context) error

	// CheckRequestAuth tries to authorize the
	// request. On siccess, the authorized user
	// object will be collected from the database
	// and set as "user" key to the request Context.
	// Otherwise, a 401 Unauthorized response will
	// be sent back.
	CheckRequestAuth(ctx *routing.Context) error
}
