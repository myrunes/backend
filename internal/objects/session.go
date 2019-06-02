package objects

import "github.com/bwmarrin/snowflake"

type Session struct {
	Key string       `json:"key"`
	UID snowflake.ID `json:"uid"`
}
