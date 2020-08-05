package caching

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
)

var (
	// expireDef is the default expiration time
	// used for cache values
	expireDef = 1 * time.Hour
)

// CacheMiddleware describes a caching module providing
// functionality to store and fetch data to/from
// a cache storage.
// This module will be initialized with SetDatabase.
// The set database must be used to fetch data from
// when a value is not found in cache. This value
// must then be saved in the cache storage.
type CacheMiddleware interface {

	// SetDatabase sets the passed database module
	// as cache storage fallback.
	SetDatabase(db database.Middleware)

	// GetUserByID returns a User object by ID
	GetUserByID(id snowflake.ID) (*objects.User, error)
	// SetUserByID sets a User object to the passed ID
	SetUserByID(id snowflake.ID, user *objects.User) error
	// GetUserByToken returns a User object by token string
	GetUserByToken(token string) (*objects.User, bool)
	// SetUserByToken sets a User object to the passed
	// token string
	SetUserByToken(token string, user *objects.User) error

	// GetPageByID returns a Page object by ID
	GetPageByID(id snowflake.ID) (*objects.Page, error)
	// SetPageByID sets a Page object to the passed ID
	SetPageByID(id snowflake.ID, page *objects.Page) error
}
