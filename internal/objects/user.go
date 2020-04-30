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

var userIDCLuster, _ = snowflake.NewNode(static.ClusterIDUsers)

var allowedUNameChars = regexp.MustCompile(`[\w_\-]+`)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	UID            snowflake.ID              `json:"uid"`
	Username       string                    `json:"username"`
	MailAddress    string                    `json:"mailaddress"`
	DisplayName    string                    `json:"displayname"`
	PassHash       []byte                    `json:"passhash,omitempty"`
	LastLogin      time.Time                 `json:"lastlogin"`
	Created        time.Time                 `json:"created"`
	Favorites      []string                  `json:"favorites"`
	PageOrder      map[string][]snowflake.ID `json:"pageorder"`
	HasOldPassword bool                      `json:"hasoldpw,omitempty"`
}

func NewUser(username, password string, authMiddleware auth.Middleware) (*User, error) {
	now := time.Now()
	passHash, err := authMiddleware.CreateHash(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Created:     now,
		LastLogin:   now,
		PassHash:    []byte(passHash),
		UID:         userIDCLuster.Generate(),
		Username:    strings.ToLower(username),
		DisplayName: username,
		Favorites:   []string{},
	}

	return user, nil
}

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

func (u *User) Validate(acceptEmptyUsername bool) error {
	if (!acceptEmptyUsername && len(u.Username) < 3) ||
		len(allowedUNameChars.FindAllString(u.Username, -1)) > 1 {

		return ErrInvalidUsername
	}

	return nil
}
