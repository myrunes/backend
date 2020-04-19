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
	keyUserByID  = "USER:ID"
	keyUserByJWT = "USER:JWT"
)

type RedisConfig struct {
	Enabled bool `json:"enabled"`

	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Redis struct {
	db database.Middleware

	client *redis.Client
}

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

	c.set(key, user, expireDef)
	return nil
}

func (c *Redis) GetUserByJWT(rawJWT string) (*objects.User, bool) {
	key := fmt.Sprintf("%s:%s", keyUserByID, rawJWT)

	var user *objects.User
	err := c.get(key, user)

	return user, err == nil && user != nil
}

func (c *Redis) SetUserByJWT(rawJWT string, user *objects.User) error {
	key := fmt.Sprintf("%s:%s", keyUserByID, rawJWT)

	return c.set(key, user, expireDef)
}

func (c *Redis) set(key string, v interface{}, expiration time.Duration) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.client.Set(key, d, expiration).Err()
}

func (c *Redis) get(key string, v interface{}) error {
	b, err := c.client.Get(key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
