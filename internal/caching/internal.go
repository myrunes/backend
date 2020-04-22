package caching

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
	"github.com/zekroTJA/timedmap"
)

const (
	secUsers = iota
	secPages
)

type Internal struct {
	db database.Middleware

	m     *timedmap.TimedMap
	users timedmap.Section
	pages timedmap.Section
}

func NewInternal() *Internal {
	tm := timedmap.New(15 * time.Minute)
	return &Internal{
		m:     tm,
		users: tm.Section(secUsers),
		pages: tm.Section(secPages),
	}
}

func (c *Internal) SetDatabase(db database.Middleware) {
	c.db = db
}

func (c *Internal) GetUserByID(id snowflake.ID) (*objects.User, error) {
	var err error
	user, ok := c.users.GetValue(id).(*objects.User)
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
		c.users.Remove(id)
	} else {
		c.users.Set(id, user, expireDef)
	}
	return nil
}

func (c *Internal) GetUserByToken(token string) (*objects.User, bool) {
	val, ok := c.users.GetValue(token).(*objects.User)
	return val, ok && val != nil
}

func (c *Internal) SetUserByToken(token string, user *objects.User) error {
	if user == nil {
		c.users.Remove(token)
	} else {
		c.users.Set(token, user, expireDef)
	}
	return nil
}

func (c *Internal) GetPageByID(id snowflake.ID) (*objects.Page, error) {
	var err error
	page, ok := c.pages.GetValue(id).(*objects.Page)
	if !ok || page == nil {
		page, err = c.db.GetPage(id)
		if err != nil {
			return nil, err
		}
		c.SetPageByID(id, page)
	}

	return page, nil
}

func (c *Internal) SetPageByID(id snowflake.ID, page *objects.Page) error {
	if page == nil {
		c.pages.Remove(id)
	} else {
		c.pages.Set(id, page, expireDef)
	}
	return nil
}
