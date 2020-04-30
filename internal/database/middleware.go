package database

import (
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/objects"
)

var (
	ErrUsernameTaken = errors.New("username taken")
)

// Middleware describes the structure of a
// database provider module.
type Middleware interface {
	// Connect  to the database using the
	// defined parameters.
	Connect(params interface{}) error
	// Close the connection to the database.
	Close()

	// CreateUser creates a new user object
	// in the database from the given user
	// object.
	CreateUser(user *objects.User) error
	// GetUser returns a user object by the
	// passed uid or username or e-mail, which
	// is passed as username parameter.
	// Therefore, the priority of matching is:
	// 1. UID, 2. username, 3. e-mail
	GetUser(uid snowflake.ID, username string) (*objects.User, error)
	// EditUser updates mutable user data values
	// of the given user in the database.
	// The passed user object may be incomplete
	// and only given and changed values must be
	// updated in the database user object.
	// If login is passed as true, the LastLogin
	// value of the user must be set to the current
	// datetime.
	EditUser(user *objects.User, login bool) (bool, error)
	// DeleteUser removes a user from the database
	// or marks it as removed so that the object
	// can not be fetched anymore.
	DeleteUser(uid snowflake.ID) error

	// CreatePage creates a page object in the
	// database from the passed page object.
	CreatePage(page *objects.Page) error
	// GetPages returns a collection of pages
	// owned by the given users uid.
	// If champion is not empty or "general",
	// only pages which champions collections
	// contain the given champion must be
	// returned.
	// If filter is not empty, only pages
	// which titles contain the filter string
	// or which champions collections contain
	// a champion which contains the filter
	// stirng must be returned.
	// If sortLess is not null, the result
	// collection must be lesss-sorted by
	// the given sortLess function.
	GetPages(uid snowflake.ID, champion, filter string, sortLess func(i, j *objects.Page) bool) ([]*objects.Page, error)
	// GetPage returns a page object by the
	// given pages uid.
	GetPage(uid snowflake.ID) (*objects.Page, error)
	EditPage(page *objects.Page) (*objects.Page, error)
	DeletePage(uid snowflake.ID) error

	SetAPIToken(token *objects.APIToken) error
	GetAPIToken(uID snowflake.ID) (*objects.APIToken, error)
	ResetAPIToken(uID snowflake.ID) error
	VerifyAPIToken(tokenStr string) (*objects.User, error)

	SetShare(share *objects.SharePage) error
	GetShare(ident string, uid, pageID snowflake.ID) (*objects.SharePage, error)
	DeleteShare(ident string, uid, pageID snowflake.ID) error
}
