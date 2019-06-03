package database

import (
	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/lol-runes/internal/objects"
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

	CreatePage(page *objects.Page) error
	GetPages(uid snowflake.ID) ([]*objects.Page, error)
	GetPage(uid snowflake.ID) (*objects.Page, error)
	EditPage(page *objects.Page) (*objects.Page, error)
	DeletePage(uid snowflake.ID) error

	CreateSession(key string, uID snowflake.ID) error
	GetSession(key string) (*objects.User, error)
}
