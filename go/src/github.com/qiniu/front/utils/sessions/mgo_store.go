package sessions

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"hd.qiniu.com/utils/object"
)

const (
	sidField       = "sid"
	valuesField    = "values"
	createdAtField = "created_at"
	updatedAtField = "updated_at"
)

type mgoStore struct {
	Id        bson.ObjectId          `bson:"_id"`
	Sid       string                 `bson:"sid"`
	Values    map[string]interface{} `bson:"values"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
}

func newMgoStore(sid string, params ...map[string]interface{}) *mgoStore {
	store := &mgoStore{
		Id:        bson.NewObjectId(),
		Sid:       sid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if len(params) > 0 {
		store.Values = params[0]
	}

	if store.Values == nil {
		store.Values = make(map[string]interface{})
	}
	return store
}

type mgoData struct {
	store    *mgoStore
	provider *MgoProvider

	changed bool // flag when values has changed
	closed  bool // flag when data destroy
}

var _ SessionStore = new(mgoData)

func newMgoData(p *MgoProvider, store *mgoStore, changed bool) *mgoData {
	sess := &mgoData{
		store:    store,
		provider: p,
		changed:  changed,
	}

	if store.Values == nil {
		store.Values = make(map[string]interface{})
		sess.changed = true
	}

	return sess
}

// get session id
func (m *mgoData) Sid() string {
	return m.store.Sid
}

func (m *mgoData) Set(key string, value interface{}) {
	m.changed = true
	m.store.Values[key] = value
}

func (m *mgoData) Get(key string) *object.Value {
	if v, ok := m.store.Values[key]; ok {
		return object.ValueTo(v)
	}
	return object.ValueTo(nil)
}

func (m *mgoData) Delete(key string) {
	m.changed = true
	delete(m.store.Values, key)
}

func (m *mgoData) Has(key string) bool {
	_, ok := m.store.Values[key]
	return ok
}

// duplicate all values
func (m *mgoData) Values() map[string]interface{} {
	values := make(map[string]interface{}, len(m.store.Values))
	for key, value := range m.store.Values {
		values[key] = value
	}
	return values
}

// clear all values in session
func (m *mgoData) Clean() {
	m.changed = true
	m.store.Values = make(map[string]interface{})
}

// save session values to store
func (m *mgoData) Flush() error {
	// has destroy
	if m.closed {
		return nil
	}

	// no changes
	if !m.changed {
		return nil
	}

	m.changed = false
	return m.provider.save(m.store)
}

// destory session values in store
func (m *mgoData) Destroy() error {
	m.closed = true
	m.store.Values = make(map[string]interface{})
	return m.provider.Destroy(m.store.Sid)
}

type MgoProvider struct {
	config  *Config
	connect func(func(c *mgo.Collection) error) error
}

var _ SessionProvider = new(MgoProvider)

func NewMgoProvider(config Config, connect func(func(c *mgo.Collection) error) error) (sess SessionProvider) {
	provider := new(MgoProvider)
	provider.config = &config
	provider.connect = connect

	connect(func(c *mgo.Collection) (err error) {

		err = c.Database.Session.Ping()
		if err != nil {
			config.Logger.Error("session mgo provider ping:", err)
			return
		}

		// sid 为唯一索引
		if err = c.EnsureIndex(mgo.Index{Key: []string{sidField}, Name: sidField, Unique: true}); err != nil {
			config.Logger.Error("session mgo provider unique EnsureIndex:", err)
			return
		}

		if config.AutoExpire {
			// 创建 mgo 自动过期的索引
			index := mgo.Index{
				Name:        updatedAtField,
				Key:         []string{updatedAtField},
				ExpireAfter: time.Duration(config.SessionExpire) * time.Second,
			}
			if err = c.EnsureIndex(index); err != nil {
				config.Logger.Error("session mgo provider expire_after EnsureIndex:", err)
				return
			}
		}
		return
	})
	return provider
}

func (p *MgoProvider) Create(sid string, params ...map[string]interface{}) (sess SessionStore, er error) {
	er = p.connect(func(c *mgo.Collection) (err error) {
		store := newMgoStore(sid, params...)

		err = c.Insert(store)

		if mgo.IsDup(err) {
			err = ErrDuplicateSid
		}

		if err == nil {
			sess = newMgoData(p, store, len(params) > 0)
		}
		return err
	})
	return
}

func (p *MgoProvider) Read(sid string) (sess SessionStore, er error) {
	er = p.connect(func(c *mgo.Collection) (err error) {
		store := &mgoStore{}
		err = c.Find(bson.M{
			sidField: sid,
		}).One(store)

		if err == mgo.ErrNotFound {
			err = ErrNotFoundSession
			return
		}

		if err == nil {
			if !p.config.AutoExpire {
				// session has expired
				if time.Since(store.UpdatedAt) > p.config.SessionExpireSeconds() {
					err = ErrNotFoundSession

					// error can secure skip
					_ = p.Destroy(sid)
					return
				}
			}

			sess = newMgoData(p, store, false)
		}
		return
	})
	return
}

// 为现有的session数据更换sid
func (p *MgoProvider) Regenerate(old string, sid string) (sess SessionStore, er error) {
	er = p.connect(func(c *mgo.Collection) (err error) {
		// 直接做更新操作
		err = c.Update(bson.M{
			sidField: old,
		}, bson.M{
			"$set": bson.M{
				sidField:       sid,
				updatedAtField: time.Now(),
			},
		})

		// 遇到重复sid
		if mgo.IsDup(err) {
			err = ErrDuplicateSid
		}

		// sid未找到
		if err == mgo.ErrNotFound {
			err = ErrNotFoundSession
		}
		return
	})

	if er != nil {
		return nil, er
	}

	return p.Read(sid)
}

func (p *MgoProvider) Destroy(sid string) (er error) {
	er = p.connect(func(c *mgo.Collection) (err error) {
		err = c.Remove(bson.M{
			sidField: sid,
		})
		return err
	})
	return
}

func (p *MgoProvider) GC() (err error) {
	if p.config.AutoExpire {
		return
	}

	err = p.connect(func(c *mgo.Collection) error {
		_, er := c.RemoveAll(bson.M{
			updatedAtField: bson.M{
				"$lte": time.Now().Add(-p.config.SessionExpireSeconds()),
			},
		})
		return er
	})
	return
}

func (p *MgoProvider) Config() *Config {
	config := *(p.config)
	return &config
}

func (p *MgoProvider) save(store *mgoStore) (err error) {
	err = p.connect(func(c *mgo.Collection) error {
		_, er := c.Upsert(bson.M{
			sidField: store.Sid,
		}, bson.M{
			sidField:       store.Sid,
			valuesField:    store.Values,
			createdAtField: store.CreatedAt,
			updatedAtField: time.Now(),
		})
		return er
	})
	return
}
