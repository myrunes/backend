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
	GetUserByJWT(rawJWT string) (*objects.User, bool)
	SetUserByJWT(rawJWT string, user *objects.User) error
}
