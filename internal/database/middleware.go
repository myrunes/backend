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

	CreateSession(key string, uID snowflake.ID) error
	GetSession(key string) (*objects.User, error)
}
