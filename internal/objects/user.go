package objects

import (
	"strings"
	"time"

	"github.com/zekroTJA/lol-runes/internal/auth"

	"github.com/bwmarrin/snowflake"
)

var userIDCLuster, _ = snowflake.NewNode(100)

type User struct {
	UID         snowflake.ID `json:"uid"`
	Username    string       `json:"username"`
	DisplayName string       `json:"displayname"`
	PassHash    []byte       `json:"passhash,omitempty"`
	LastLogin   time.Time    `json:"lastlogin"`
	Created     time.Time    `json:"created"`
}

func NewUser(username, password string, authMiddleware auth.Middleware) (*User, error) {
	now := time.Now()
	passHash, err := authMiddleware.CreateHash([]byte(password))
	if err != nil {
		return nil, err
	}

	user := &User{
		Created:     now,
		LastLogin:   now,
		PassHash:    passHash,
		UID:         userIDCLuster.Generate(),
		Username:    strings.ToLower(username),
		DisplayName: username,
	}

	return user, nil
}
