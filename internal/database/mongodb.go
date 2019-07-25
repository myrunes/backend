package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/zekroTJA/myrunes/internal/objects"

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
	sessions,
	apitokens,
	shares *mongo.Collection
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
		users:     m.db.Collection("users"),
		pages:     m.db.Collection("pages"),
		sessions:  m.db.Collection("sessions"),
		shares:    m.db.Collection("shares"),
		apitokens: m.db.Collection("apitokens"),
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

func (m *MongoDB) EditUser(user *objects.User, login bool) (bool, error) {
	oldUser, err := m.GetUser(user.UID, "")
	if err != nil {
		return false, err
	}

	if oldUser == nil {
		return false, nil
	}

	if login {
		oldUser.LastLogin = time.Now()
	}

	if user.DisplayName != "" {
		oldUser.DisplayName = user.DisplayName
	}

	if user.Favorites != nil {
		oldUser.Favorites = user.Favorites
	}

	if user.Username != "" {
		u, err := m.GetUser(snowflake.ID(-1), user.Username)
		if err != nil {
			return false, err
		}
		if u != nil && u.UID != oldUser.UID {
			return false, ErrUsernameTaken
		}
		oldUser.Username = user.Username
	}

	if user.PassHash != nil && len(user.PassHash) > 0 {
		oldUser.PassHash = user.PassHash
	}

	return true, m.insertOrUpdate(m.collections.users,
		bson.M{"uid": oldUser.UID}, oldUser)
}

func (m *MongoDB) DeleteUser(uid snowflake.ID) error {
	_, err := m.collections.users.DeleteOne(ctxTimeout(5*time.Second), bson.M{"uid": uid})
	if err != nil {
		return err
	}

	_, err = m.collections.pages.DeleteMany(ctxTimeout(5*time.Second),
		bson.M{"owner": uid})

	return err
}

func (m *MongoDB) CreatePage(page *objects.Page) error {
	return m.insert(m.collections.pages, page)
}

func (m *MongoDB) GetPages(uid snowflake.ID) ([]*objects.Page, error) {
	count, err := m.count(m.collections.pages, bson.M{"owner": uid})
	if err != nil {
		return nil, err
	}

	pages := make([]*objects.Page, count)

	if count == 0 {
		return pages, nil
	}

	res, err := m.collections.pages.Find(ctxTimeout(5*time.Second), bson.M{"owner": uid})
	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next(ctxTimeout(2 * time.Second)) {
		page := new(objects.Page)
		err = res.Decode(page)
		if err != nil {
			return nil, err
		}
		pages[i] = page
		i++
	}

	return pages, nil
}

func (m *MongoDB) GetPage(uid snowflake.ID) (*objects.Page, error) {
	page := new(objects.Page)
	ok, err := m.get(m.collections.pages, bson.M{"uid": uid}, page)
	if err != nil || !ok {
		return nil, err
	}
	return page, nil
}

func (m *MongoDB) EditPage(page *objects.Page) (*objects.Page, error) {
	oldPage, err := m.GetPage(page.UID)
	if err != nil {
		return nil, err
	}
	if oldPage == nil {
		return nil, nil
	}

	page.Created = oldPage.Created
	page.UID = oldPage.UID
	page.Owner = oldPage.Owner
	page.Edited = time.Now()
	err = page.Validate()
	if err != nil {
		return nil, err
	}

	return page, m.insertOrUpdate(m.collections.pages, bson.M{"uid": page.UID}, page)
}

func (m *MongoDB) DeletePage(uid snowflake.ID) error {
	_, err := m.collections.pages.DeleteOne(ctxTimeout(5*time.Second), bson.M{"uid": uid})
	return err
}

func (m *MongoDB) CreateSession(key string, uID snowflake.ID, expires time.Time, addr string) error {
	session := objects.NewSession(key, uID, expires, addr)

	return m.insert(m.collections.sessions, session)
}

func (m *MongoDB) GetSession(key string, addr string) (*objects.User, error) {
	session := new(objects.Session)
	ok, err := m.get(m.collections.sessions, bson.M{"key": key}, session)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	if time.Now().After(session.Expires) {
		return nil, m.DeleteSession("", session.SessionID)
	}

	user := new(objects.User)
	ok, err = m.get(m.collections.users, bson.M{"uid": session.UID}, user)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	session.LastAccess = time.Now()
	session.LastAccessIP = addr
	err = m.insertOrUpdate(m.collections.sessions, bson.M{"sessionid": session.SessionID}, session)

	return user, err
}

func (m *MongoDB) GetSessions(uID snowflake.ID) ([]*objects.Session, error) {
	count, err := m.count(m.collections.sessions, bson.M{"uid": uID})
	if err != nil {
		return nil, err
	}

	sessions := make([]*objects.Session, count)

	if count == 0 {
		return sessions, nil
	}

	res, err := m.collections.sessions.Find(ctxTimeout(5*time.Second), bson.M{"uid": uID})
	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next(ctxTimeout(5 * time.Second)) {
		if int64(i) >= count {
			break
		}

		s := new(objects.Session)
		if err = res.Decode(s); err != nil {
			return nil, err
		}
		if time.Now().Before(s.Expires) {
			sessions[i] = s
			i++
		} else {
			// m.DeleteSession("", s.SessionID)
		}
	}

	return sessions[:i], nil
}

func (m *MongoDB) DeleteSession(key string, sessionID snowflake.ID) error {
	_, err := m.collections.sessions.DeleteOne(ctxTimeout(5*time.Second),
		bson.M{"$or": bson.A{
			bson.M{"key": key},
			bson.M{"sessionid": sessionID},
		}})
	return err
}

func (m *MongoDB) CleanupExpiredSessions() error {
	_, err := m.collections.sessions.DeleteMany(ctxTimeout(5*time.Second),
		bson.M{
			"expires": bson.M{
				"$lte": time.Now(),
			},
		})

	return err
}

func (m *MongoDB) SetAPIToken(token *objects.APIToken) error {
	return m.insertOrUpdate(m.collections.apitokens, &bson.M{"userid": token.UserID}, token)
}

func (m *MongoDB) GetAPIToken(uID snowflake.ID) (*objects.APIToken, error) {
	token := new(objects.APIToken)
	ok, err := m.get(m.collections.apitokens, bson.M{"userid": uID}, token)
	if err != nil || !ok {
		return nil, err
	}
	return token, nil
}

func (m *MongoDB) ResetAPIToken(uID snowflake.ID) error {
	_, err := m.collections.apitokens.DeleteOne(ctxTimeout(5*time.Second), bson.M{"userid": uID})
	return err
}

func (m *MongoDB) VerifyAPIToken(tokenStr string) (*objects.User, error) {
	token := new(objects.APIToken)
	ok, err := m.get(m.collections.apitokens, bson.M{"token": tokenStr}, token)
	if err != nil || !ok {
		return nil, err
	}

	return m.GetUser(token.UserID, "")
}

func (m *MongoDB) SetShare(share *objects.SharePage) error {
	return m.insertOrUpdate(m.collections.shares, bson.M{
		"$or": bson.A{
			bson.M{"uid": share.UID},
			bson.M{"pageid": share.PageID},
		},
	}, share)
}

func (m *MongoDB) GetShare(ident string, uid, pageID snowflake.ID) (*objects.SharePage, error) {
	share := new(objects.SharePage)

	ok, err := m.get(m.collections.shares, bson.M{
		"$or": bson.A{
			bson.M{"ident": ident},
			bson.M{"uid": uid},
			bson.M{"pageid": pageID},
		},
	}, share)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return share, nil
}

func (m *MongoDB) DeleteShare(ident string, uid, pageID snowflake.ID) error {
	_, err := m.collections.shares.DeleteOne(ctxTimeout(5*time.Second), bson.M{
		"$or": bson.A{
			bson.M{"ident": ident},
			bson.M{"uid": uid},
			bson.M{"pageid": pageID},
		},
	})

	return err
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

	if res.MatchedCount == 0 {
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

func (M *MongoDB) count(collection *mongo.Collection, filter interface{}) (int64, error) {
	return collection.CountDocuments(ctxTimeout(5*time.Second), filter)
}

func ctxTimeout(d time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), d)
	return ctx
}
