package database

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/myrunes/backend/internal/objects"

	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

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

	ctxConnect, cancelConnect := ctxTimeout(5 * time.Second)
	defer cancelConnect()

	if err = m.client.Connect(ctxConnect); err != nil {
		return
	}

	ctxPing, cancelPing := ctxTimeout(5 * time.Second)
	defer cancelPing()

	if err = m.client.Ping(ctxPing, readpref.Primary()); err != nil {
		return err
	}

	m.db = m.client.Database(cfg.DataDB)

	m.collections = &collections{
		users:     m.db.Collection("users"),
		pages:     m.db.Collection("pages"),
		shares:    m.db.Collection("shares"),
		apitokens: m.db.Collection("apitokens"),
	}

	return err
}

func (m *MongoDB) Close() {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	m.client.Disconnect(ctx)
}

func (m *MongoDB) CreateUser(user *objects.User) error {
	return m.insert(m.collections.users, user)
}

func (m *MongoDB) GetUser(uid snowflake.ID, username string) (*objects.User, error) {
	user := new(objects.User)

	ok, err := m.get(m.collections.users, bson.M{"$or": bson.A{
		bson.M{"username": username},
		bson.M{"mailaddress": username},
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

	if user.PageOrder != nil {
		oldUser.PageOrder = user.PageOrder
	}

	if user.MailAddress != "" {
		if user.MailAddress == "__RESET__" {
			oldUser.MailAddress = ""
		} else {
			oldUser.MailAddress = user.MailAddress
		}
	}

	return true, m.insertOrUpdate(m.collections.users,
		bson.M{"uid": oldUser.UID}, oldUser)
}

func (m *MongoDB) DeleteUser(uid snowflake.ID) error {
	ctxDelOne, cancelDelOne := ctxTimeout(5 * time.Second)
	defer cancelDelOne()

	_, err := m.collections.users.DeleteOne(ctxDelOne, bson.M{"uid": uid})
	if err != nil {
		return err
	}

	ctxDelMany, cancelDelMany := ctxTimeout(5 * time.Second)
	defer cancelDelMany()

	_, err = m.collections.pages.DeleteMany(ctxDelMany,
		bson.M{"owner": uid})

	return err
}

func (m *MongoDB) CreatePage(page *objects.Page) error {
	return m.insert(m.collections.pages, page)
}

func (m *MongoDB) GetPages(uid snowflake.ID, champion, filter string, sortLess func(i, j *objects.Page) bool) ([]*objects.Page, error) {
	var query bson.M
	if champion != "" && champion != "general" {
		query = bson.M{"owner": uid, "champions": champion}
	} else {
		query = bson.M{"owner": uid}
	}

	if filter != "" {
		query["$or"] = bson.A{
			bson.M{
				"title": bson.M{
					"$regex": fmt.Sprintf("(?i).*%s.*", filter),
				},
			},
			bson.M{
				"champions": bson.M{
					"$regex": fmt.Sprintf("(?i).*%s.*", filter),
				},
			},
		}
	}

	count, err := m.count(m.collections.pages, query)
	if err != nil {
		return nil, err
	}

	pages := make([]*objects.Page, count)

	if count == 0 {
		return pages, nil
	}

	ctxFind, cancelFind := ctxTimeout(5 * time.Second)
	defer cancelFind()

	res, err := m.collections.pages.Find(ctxFind, query)
	if err != nil {
		return nil, err
	}

	ctxNext, cancelNext := ctxTimeout(5 * time.Second)
	defer cancelNext()

	i := 0
	for res.Next(ctxNext) {
		page := new(objects.Page)
		err = res.Decode(page)
		if err != nil {
			return nil, err
		}
		pages[i] = page
		i++
	}

	if sortLess != nil {
		sort.Slice(pages, func(i, j int) bool {
			return sortLess(pages[i], pages[j])
		})
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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err := m.collections.pages.DeleteOne(ctx, bson.M{"uid": uid})
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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err := m.collections.apitokens.DeleteOne(ctx, bson.M{"userid": uID})
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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err := m.collections.shares.DeleteOne(ctx, bson.M{
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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, v)
	return err
}

func (m *MongoDB) insertOrUpdate(collection *mongo.Collection, filter, obj interface{}) error {
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

	res, err := collection.UpdateOne(
		ctx,
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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()

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
	ctx, cancel := ctxTimeout(5 * time.Second)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}

func ctxTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	return ctx, cancel
}
