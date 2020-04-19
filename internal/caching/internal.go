package caching

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
	"github.com/zekroTJA/timedmap"
)

type Internal struct {
	db database.Middleware

	m *timedmap.TimedMap
}

func NewInternal() *Internal {
	return &Internal{
		m: timedmap.New(15 * time.Minute),
	}
}

func (c *Internal) SetDatabase(db database.Middleware) {
	c.db = db
}

func (c *Internal) GetUserByID(id snowflake.ID) (*objects.User, error) {
	user, ok := c.m.GetValue(id).(*objects.User)
	if !ok {
		user, err := c.db.GetUser(id, "")
		if err != nil {
			return nil, err
		}
		c.m.Set(id, user, expireDef)
	}

	return user, nil
}

func (c *Internal) GetUserByJWT(rawJWT string) (*objects.User, bool) {
	val, ok := c.m.GetValue(rawJWT).(*objects.User)
	return val, ok
}

func (c *Internal) SetUserByJWT(rawJWT string, user *objects.User) error {
	c.m.Set(rawJWT, user, expireDef)
	return nil
}
