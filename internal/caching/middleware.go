package caching

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
)

var (
	expireDef = 1 * time.Hour
)

type Middleware interface {
	SetDatabase(db database.Middleware)

	GetUserByID(id snowflake.ID) (*objects.User, error)
	SetUserByID(id snowflake.ID, user *objects.User) error
	GetUserByToken(token string) (*objects.User, bool)
	SetUserByToken(token string, user *objects.User) error

	GetPageByID(id snowflake.ID) (*objects.Page, error)
	SetPageByID(id snowflake.ID, page *objects.Page) error
}
