package objects

import (
	"time"

	"github.com/myrunes/backend/pkg/random"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/static"
)

// shareIDNode is the node to generate share snowflake IDs.
var shareIDNode, _ = snowflake.NewNode(static.NodeIDShares)

// SharePage wraps a RunePage public share.
type SharePage struct {
	UID         snowflake.ID `json:"uid"`
	Ident       string       `json:"ident"`
	OwnerID     snowflake.ID `json:"owner"`
	PageID      snowflake.ID `json:"page"`
	Created     time.Time    `json:"created"`
	MaxAccesses int          `json:"maxaccesses"`
	Expires     time.Time    `json:"expires"`
	Accesses    int          `json:"accesses"`
	LastAccess  time.Time    `json:"lastaccess"`
	AccessIPs   []string     `json:"accessips,omitempty"`
}

// NEwSharePage creates a new SharePage instance with
// the passed ownerID, pageID, maxAccess count and
// expiration time.
// If maxAccesses is 0, maxAccesses is set to -1 which
// indicates that the share is not access limited by
// access count.
// If expire time is the default TIme object Time{},
// the expiration will be set to 100 years, which
// should be enough to count as unlimited access
// by time.
func NewSharePage(ownerID, pageID snowflake.ID, maxAccesses int, expires time.Time) (*SharePage, error) {
	now := time.Now()
	var err error

	if (expires == time.Time{}) {
		expires = now.Add(100 * 365 * 24 * time.Hour)
	}

	if maxAccesses == 0 {
		maxAccesses = -1
	}

	share := &SharePage{
		Accesses:    0,
		Created:     now,
		Expires:     expires,
		LastAccess:  now,
		MaxAccesses: maxAccesses,
		OwnerID:     ownerID,
		PageID:      pageID,
		UID:         shareIDNode.Generate(),
		AccessIPs:   make([]string, 0),
	}

	const identSubset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	share.Ident, err = random.String(5, identSubset)

	return share, err
}
