package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/zekroTJA/lol-runes/internal/objects"

	"github.com/bwmarrin/snowflake"
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
	sessions *mongo.Collection
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

	if err = m.client.Ping(ctxTimeout(2*time.Second), readpref.Primary()); err != nil {
		return err
	}

	m.db = m.client.Database(cfg.DataDB)

	m.collections = &collections{
		users:    m.db.Collection("users"),
		pages:    m.db.Collection("pages"),
		sessions: m.db.Collection("sessions"),
	}

	return err
}

func (m *MongoDB) Close() {
	m.client.Disconnect(ctxTimeout(5 * time.Second))
}

func (m *MongoDB) CreateUser(user *objects.User) error {
	return m.insert(m.collections.users, user)
}

func (m *MongoDB) GetUser(uid snowflake.ID, username string) (*objects.User, error) {
	user := new(objects.User)

	ok, err := m.get(m.collections.users, bson.M{"$or": bson.A{
		bson.M{"username": username},
		bson.M{"uid": uid},
	}}, user)

	if !ok {
		user = nil
	}

	return user, err
}

func (m *MongoDB) CreateSession(key string, uID snowflake.ID) error {
	session := &objects.Session{
		Key: key,
		UID: uID,
	}

	return m.insertOrUpdate(m.collections.sessions, bson.M{"uid": uID}, session)
}

func (m *MongoDB) GetSession(key string) (*objects.User, error) {
	session := new(objects.Session)
	ok, err := m.get(m.collections.sessions, bson.M{"key": key}, session)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	user := new(objects.User)
	ok, err = m.get(m.collections.users, bson.M{"uid": session.UID}, user)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return user, nil
}

// --- HELPERS ------------------------------------------------------------------

func (m *MongoDB) insert(collection *mongo.Collection, v interface{}) error {
	_, err := collection.InsertOne(ctxTimeout(5*time.Second), v)
	return err
}

func (m *MongoDB) insertOrUpdate(collection *mongo.Collection, filter, obj interface{}) error {
	res, err := collection.UpdateOne(
		ctxTimeout(5*time.Second),
		filter, bson.M{
			"$set": obj,
		})

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return m.insert(collection, obj)
	}

	return err
}

func (m *MongoDB) get(collection *mongo.Collection, filter interface{}, v interface{}) (bool, error) {
	ctx := ctxTimeout(5 * time.Second)
	res := collection.FindOne(ctx, filter)

	if res == nil {
		return false, nil
	}

	err := res.Decode(v)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	return true, nil
}

func ctxTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
