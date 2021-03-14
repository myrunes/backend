package database

import (
	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/objects"
)

// Middleware describes the structure of a
// database provider module.
//
// On fetching values from the database
// which dont exist, an error returning
// is not expected. Errors should only be
// returned if something went wrong
// accessing the database and will be
// returned as 500 Internal Error from
// the REST API. If an non-existing
// object was fetched, only return nil
// (or the default value for the type)
// for both the value and the error.
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
	// EditUser updates a user object in the
	// database to the object passed by its
	// UID.
	EditUser(user *objects.User) error
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
	GetPages(
		uid snowflake.ID,
		champion,
		filter string,
		sortLess func(i, j *objects.Page) bool,
	) ([]*objects.Page, error)
	// GetPage returns a page object by the
	// given pages uid.
	GetPage(uid snowflake.ID) (*objects.Page, error)
	// EditPage replaces the page object in
	// the database by the passed page
	// object by its UID.
	EditPage(page *objects.Page) error
	// DeletePage removes a page object from
	// the database or marks it as removed
	// so it's not accessable anymore.
	DeletePage(uid snowflake.ID) error
	// DeleteUserPages deletes all pages
	// of the users UID passed.
	DeleteUserPages(uid snowflake.ID) error

	// GetRefreshToken returns a refresh token object
	// from the database matching the given refresh
	// token string.
	GetRefreshToken(token string) (*objects.RefreshToken, error)
	// GetRefreshTokens returns a list of refresh tokens
	// belonging to the given userID.
	GetRefreshTokens(userID snowflake.ID) ([]*objects.RefreshToken, error)
	// SetRefreshToken sets a given refresh token
	// object to the database or updates one.
	SetRefreshToken(t *objects.RefreshToken) error
	// RemoveRefreshToken removes a refresh token from
	// database if existent by the given token.
	RemoveRefreshToken(id snowflake.ID) error
	// CleanupExpiredTokens removes all expired tokens
	// from the database.
	CleanupExpiredTokens() (int, error)

	// SetAPIToken sets the passed API token
	// to the user defined in the APIToken
	// object.
	SetAPIToken(token *objects.APIToken) error
	// GetAPIToken returns the APIToken object,
	// if available, of the passed users uid.
	GetAPIToken(uid snowflake.ID) (*objects.APIToken, error)
	// ResetAPIToken deletes the APIToken
	// object of the passed users uid so
	// that it is no more accessable.
	ResetAPIToken(uid snowflake.ID) error
	// VerifyAPIToken returns a User object
	// which the passed API token string
	// belongs to.
	VerifyAPIToken(tokenStr string) (*objects.User, error)

	// SetShare creates a nnew share entry
	// in the database from the passed SharePage
	// object.
	SetShare(share *objects.SharePage) error
	// GetShare returns the SharePage object by
	// the shares ident, uid or pageID of the
	// RunePage the share is assigned to.
	// (Priority in this order)
	GetShare(ident string, uid, pageID snowflake.ID) (*objects.SharePage, error)
	// DeleteShare removes a SharePage object
	// from the database or makes it inaccessable
	// by the shares ident, uid oder the pageID
	// of the RunePage the share is belonging to.
	// (Priority in this order)
	DeleteShare(ident string, uid, pageID snowflake.ID) error
}
