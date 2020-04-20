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
	var err error
	user, ok := c.m.GetValue(id).(*objects.User)
	if !ok || user == nil {
		user, err = c.db.GetUser(id, "")
		if err != nil {
			return nil, err
		}
		c.SetUserByID(id, user)
	}

	return user, nil
}

func (c *Internal) SetUserByID(id snowflake.ID, user *objects.User) error {
	if user == nil {
		c.m.Remove(id)
	} else {
		c.m.Set(id, user, expireDef)
	}
	return nil
}

func (c *Internal) GetUserByJWT(rawJWT string) (*objects.User, bool) {
	val, ok := c.m.GetValue(rawJWT).(*objects.User)
	return val, ok && val != nil
}

func (c *Internal) SetUserByJWT(rawJWT string, user *objects.User) error {
	if user == nil {
		c.m.Remove(rawJWT)
	} else {
		c.m.Set(rawJWT, user, expireDef)
	}
	return nil
}
