package caching

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis"
	"github.com/myrunes/backend/internal/database"
	"github.com/myrunes/backend/internal/objects"
)

const (
	keyUserByID    = "USER:ID"
	keyUserByToken = "USER:TK"
	keyPageByID    = "PAGE:ID"
)

// RedisConfig contains configuration
// values for the Redis Database
// connection.
type RedisConfig struct {
	Enabled bool `json:"enabled"`

	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Redis provides a caching module which
// uses Redis to store and manage cache
// values.
type Redis struct {
	db database.Middleware

	client *redis.Client
}

// NewRedis creates a new instance of
// Redis with the given RedisConfig
// instance cfg.
func NewRedis(cfg *RedisConfig) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

func (c *Redis) SetDatabase(db database.Middleware) {
	c.db = db
}

func (c *Redis) GetUserByID(id snowflake.ID) (*objects.User, error) {
	key := fmt.Sprintf("%s:%d", keyUserByID, id)

	var user *objects.User
	err := c.get(key, user)
	if err != nil || user == nil {
		user, err = c.db.GetUser(id, "")
		if err != nil {
			return nil, err
		}
		c.SetUserByID(id, user)
	}

	return user, nil
}

func (c *Redis) SetUserByID(id snowflake.ID, user *objects.User) error {
	key := fmt.Sprintf("%s:%d", keyUserByID, id)

	if user == nil {
		return c.set(key, nil, expireDef)
	}
	return c.set(key, user, expireDef)
}

func (c *Redis) GetUserByToken(token string) (*objects.User, bool) {
	key := fmt.Sprintf("%s:%s", keyUserByToken, token)

	var user *objects.User
	err := c.get(key, user)

	return user, err == nil && user != nil
}

func (c *Redis) SetUserByToken(token string, user *objects.User) error {
	key := fmt.Sprintf("%s:%s", keyUserByToken, token)

	if user == nil {
		return c.set(key, nil, expireDef)
	}
	return c.set(key, user, expireDef)
}

func (c *Redis) GetPageByID(id snowflake.ID) (*objects.Page, error) {
	key := fmt.Sprintf("%s:%d", keyPageByID, id)

	var page *objects.Page
	err := c.get(key, page)
	if err != nil || page == nil {
		page, err = c.db.GetPage(id)
		if err != nil {
			return nil, err
		}
		c.SetPageByID(id, page)
	}

	return page, nil
}

func (c *Redis) SetPageByID(id snowflake.ID, page *objects.Page) error {
	key := fmt.Sprintf("%s:%d", keyPageByID, id)

	if page == nil {
		return c.set(key, nil, expireDef)
	}
	return c.set(key, page, expireDef)
}

// set sets a value in the database to the given key with the
// defined expiration duration.
// The value v must be a reference to a JSON serializable
// object instance.
func (c *Redis) set(key string, v interface{}, expiration time.Duration) error {
	if v == nil {
		return c.client.Del(key).Err()
	}

	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.client.Set(key, d, expiration).Err()
}

// get fetches a value from the database by key and writes
// the result to v.
// The value v must be a reference to a JSON serializable
// object instance.
func (c *Redis) get(key string, v interface{}) error {
	b, err := c.client.Get(key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
