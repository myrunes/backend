package objects

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/myrunes/backend/internal/auth"
	"github.com/myrunes/backend/internal/static"

	"github.com/bwmarrin/snowflake"
)

// userIDNode is the node to generate user snowflake IDs.
var userIDNode, _ = snowflake.NewNode(static.NodeIDUsers)

// allowedUNameChars is a regular expression which matches
// on user name strings which are valid.
var allowedUNameChars = regexp.MustCompile(`[\w_\-]+`)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

// User wraps a general user object.
type User struct {
	UID            snowflake.ID              `json:"uid"`
	Username       string                    `json:"username"`
	MailAddress    string                    `json:"mailaddress,omitempty"`
	DisplayName    string                    `json:"displayname"`
	LastLogin      time.Time                 `json:"lastlogin,omitempty"`
	Created        time.Time                 `json:"created"`
	Favorites      []string                  `json:"favorites,omitempty"`
	PageOrder      map[string][]snowflake.ID `json:"pageorder,omitempty"`
	HasOldPassword bool                      `json:"hasoldpw,omitempty"`

	PassHash []byte `json:"-"`
}

// NewUser creates a new User object with the given
// username and password which will be hashed using
// the passed authModdleware and then saved to the
// user object.
func NewUser(username, password string, authMiddleware auth.AuthMiddleware) (*User, error) {
	now := time.Now()
	passHash, err := authMiddleware.CreateHash(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Created:     now,
		LastLogin:   now,
		PassHash:    []byte(passHash),
		UID:         userIDNode.Generate(),
		Username:    strings.ToLower(username),
		DisplayName: username,
		Favorites:   []string{},
	}

	return user, nil
}

// Update sets mutable user data to the
// current user object from the given
// newUser object properties.
// If login is set to true, lastLogin
// will be set to the current time.
func (u *User) Update(newUser *User, login bool) {
	if login {
		u.LastLogin = time.Now()
	}

	if newUser == nil {
		return
	}

	if newUser.DisplayName != "" {
		u.DisplayName = newUser.DisplayName
	}

	if newUser.Favorites != nil {
		u.Favorites = newUser.Favorites
	}

	if newUser.Username != "" {
		u.Username = newUser.Username
	}

	if newUser.PassHash != nil && len(newUser.PassHash) > 0 {
		u.PassHash = newUser.PassHash
	}

	if newUser.PageOrder != nil {
		u.PageOrder = newUser.PageOrder
	}

	if newUser.MailAddress != "" {
		if newUser.MailAddress == "__RESET__" {
			u.MailAddress = ""
		} else {
			u.MailAddress = newUser.MailAddress
		}
	}
}

// Validate checks if the user object
// is built by specification.
// If the validation fails, the failure
// will be returned as error object.
func (u *User) Validate(acceptEmptyUsername bool) error {
	if (!acceptEmptyUsername && len(u.Username) < 3) ||
		len(allowedUNameChars.FindAllString(u.Username, -1)) > 1 {

		return ErrInvalidUsername
	}

	return nil
}

// Sanitize creates a new User object from
// the current User object which only contains
// information which shall be publicly visible.
func (u *User) Sanitize() *User {
	return &User{
		UID:         u.UID,
		Created:     u.Created,
		DisplayName: u.DisplayName,
		Username:    u.Username,
	}
}
