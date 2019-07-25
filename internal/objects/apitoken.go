package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type APIToken struct {
	UserID  snowflake.ID `json:"userid"`
	Token   string       `json:"token"`
	Created time.Time    `json:"created"`
}
