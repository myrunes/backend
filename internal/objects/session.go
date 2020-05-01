package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/static"
)

// sessionIDNode is the node to generate session snowflake IDs.
var sessionIDNode, _ = snowflake.NewNode(static.NodeIDSessions)

// Session wraps a brwoser login session instance.
type Session struct {
	SessionID    snowflake.ID `json:"sessionid"`
	Key          string       `json:"key"`
	UID          snowflake.ID `json:"uid"`
	Expires      time.Time    `json:"expires"`
	LastAccess   time.Time    `json:"lastaccess"`
	LastAccessIP string       `json:"lastaccessip"`
}

// NewSession creates a session instance with the
// given session key, expires time, last access IP
// addr - belonging to the passed user ID (uID).
func NewSession(key string, uID snowflake.ID, expires time.Time, addr string) *Session {
	return &Session{
		Key:          key,
		UID:          uID,
		Expires:      expires,
		LastAccessIP: addr,
		LastAccess:   time.Now(),
		SessionID:    sessionIDNode.Generate(),
	}
}
