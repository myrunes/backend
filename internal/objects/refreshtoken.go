package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type RefreshToken struct {
	Token    string
	UserID   snowflake.ID
	Deadline time.Time
}

type AccessToken struct {
	Token string `json:"accesstoken"`
}
