package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

// APIToken wraps an API access token string
// with the UserID belonging to the token and
// a time when the token was created.
type APIToken struct {
	UserID  snowflake.ID `json:"userid"`
	Token   string       `json:"token"`
	Created time.Time    `json:"created"`
}
