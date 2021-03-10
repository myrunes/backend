package objects

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/static"
)

var refreshTokenIDNode, _ = snowflake.NewNode(static.NodeIDRefreshTokens)

type RefreshToken struct {
	ID               snowflake.ID `json:"id"`
	Token            string       `json:"token,omitempty"`
	UserID           snowflake.ID `json:"userid"`
	Deadline         time.Time    `json:"deadline"`
	LastAccess       time.Time    `json:"lastaccess"`
	LastAccessClient string       `json:"lastaccessclient"`
	LastAccessIP     string       `json:"lastaccessip"`
}

type AccessToken struct {
	Token string `json:"accesstoken"`
}

func (rt *RefreshToken) SetID() *RefreshToken {
	rt.ID = refreshTokenIDNode.Generate()
	return rt
}

func (rt *RefreshToken) Sanitize() {
	rt.Token = ""
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.Deadline)
}
