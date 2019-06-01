package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongodb://myDBReader:D1fficultP%40ssw0rd@mongodb0.example.com:27017/admin

type MongoDB struct {
	client      *mongo.Client
	db          *mongo.Database
	collections *collections
}

type MongoConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	AuthDB   string `json:"auth_db"`
	DataDB   string `json:"data_db"`
}

type collections struct {
	users,
	pages,
	logins *mongo.Collection
}

func (m *MongoDB) Connect(params interface{}) (err error) {
	cfg, ok := params.(*MongoConfig)
	if !ok {
		return errors.New("invalid config data type")
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.AuthDB)
	if m.client, err = mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		return
	}

	if err = m.client.Connect(ctxTimeout(10 * time.Second)); err != nil {
		return
	}

	err = m.client.Ping(ctxTimeout(2*time.Second), readpref.Primary())

	m.db = m.client.Database(cfg.AuthDB)

	m.collections = &collections{
		users:  m.db.Collection("users"),
		pages:  m.db.Collection("pages"),
		logins: m.db.Collection("logins"),
	}

	return
}

func (m *MongoDB) Close() {
	m.client.Disconnect(ctxTimeout(5 * time.Second))
}

// --- HELPERS ------------------------------------------------------------------

func (m *MongoDB) insert(collection *mongo.Collection, v interface{}) error {
	_, err := collection.InsertOne(ctxTimeout(5*time.Second), v)
	return err
}

func ctxTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
