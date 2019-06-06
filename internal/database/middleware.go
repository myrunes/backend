package database

import (
	"errors"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/myrunes/internal/objects"
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
	GetPages(uid snowflake.ID) ([]*objects.Page, error)
	GetPage(uid snowflake.ID) (*objects.Page, error)
	EditPage(page *objects.Page) (*objects.Page, error)
	DeletePage(uid snowflake.ID) error

	CreateSession(key string, uID snowflake.ID, expires time.Time) error
	GetSession(key string) (*objects.User, error)
	GetSessions(uID snowflake.ID) ([]*objects.Session, error)
	DeleteSession(key string) error
}
