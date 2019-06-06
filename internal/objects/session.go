package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var sessionIDCluster, _ = snowflake.NewNode(300)

type Session struct {
	SessionID    snowflake.ID `json:"sessionid"`
	Key          string       `json:"key"`
	UID          snowflake.ID `json:"uid"`
	Expires      time.Time    `json:"expires"`
	LastAccess   time.Time    `json:"lastaccess"`
	LastAccessIP string       `json:"lastaccessip"`
}

func NewSession(key string, uID snowflake.ID, expires time.Time, addr string) *Session {
	return &Session{
		Key:          key,
		UID:          uID,
		Expires:      expires,
		LastAccessIP: addr,
		LastAccess:   time.Now(),
		SessionID:    sessionIDCluster.Generate(),
	}
}
