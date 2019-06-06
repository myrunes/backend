package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

type Session struct {
	Key     string       `json:"key"`
	UID     snowflake.ID `json:"uid"`
	Expires time.Time    `json:"expires"`
}
