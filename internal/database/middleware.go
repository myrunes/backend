package database

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/myrunes/internal/objects"
)

var (
	ErrUsernameTaken = errors.New("username taken")
)

// Middleware describes the structure of a
// database middleware.
type Middleware interface {
	// Connect to the database server or file or
	// whatever you are about to use.
	Connect(params interface{}) error
	// Close the connection to the database.
	Close()

	CreateUser(user *objects.User) error
	GetUser(uid snowflake.ID, username string) (*objects.User, error)
	EditUser(user *objects.User, login bool) (bool, error)
	DeleteUser(uid snowflake.ID) error

	CreatePage(page *objects.Page) error
	GetPages(uid snowflake.ID, champion string, sortLess func(i, j *objects.Page) bool) ([]*objects.Page, error)
	GetPage(uid snowflake.ID) (*objects.Page, error)
	EditPage(page *objects.Page) (*objects.Page, error)
	DeletePage(uid snowflake.ID) error

	CreateSession(key string, uID snowflake.ID, expires time.Time, addr string) error
	GetSession(key string, addr string) (*objects.User, error)
	GetSessions(uID snowflake.ID) ([]*objects.Session, error)
	DeleteSession(key string, sessionID snowflake.ID) error
	CleanupExpiredSessions() error

	SetAPIToken(token *objects.APIToken) error
	GetAPIToken(uID snowflake.ID) (*objects.APIToken, error)
	ResetAPIToken(uID snowflake.ID) error
	VerifyAPIToken(tokenStr string) (*objects.User, error)

	SetShare(share *objects.SharePage) error
	GetShare(ident string, uid, pageID snowflake.ID) (*objects.SharePage, error)
	DeleteShare(ident string, uid, pageID snowflake.ID) error
}
