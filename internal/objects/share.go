package objects

import (
	"time"

	"github.com/zekroTJA/myrunes/pkg/random"

	"github.com/bwmarrin/snowflake"
	"github.com/zekroTJA/myrunes/internal/static"
)

var shareIDCLuster, _ = snowflake.NewNode(static.ClusterIDShares)

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
		UID:         shareIDCLuster.Generate(),
		AccessIPs:   make([]string, 0),
	}

	share.Ident, err = random.GetRandBase64Str(5)

	return share, err
}
