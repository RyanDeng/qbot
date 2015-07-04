package cache

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"hd.qiniu.com/utils/object"
)

type MgoConfig struct {
	AutoExpire     bool
	KeyField       string
	ValueField     string
	ExpiredAtField string
}

func NewMgoConfig() MgoConfig {
	config := MgoConfig{
		KeyField:       "key",
		ValueField:     "value",
		ExpiredAtField: "expired_at",
	}

	return config
}

type MgoCache struct {
	config  *MgoConfig
	connect func(func(c *mgo.Collection) error) error
}

func NewMgoCache(config MgoConfig, connect func(func(c *mgo.Collection) error) error) (cache Cache, err error) {
	cache = &MgoCache{
		config:  &config,
		connect: connect,
	}

	connect(func(c *mgo.Collection) error {
		err = c.Database.Session.Ping()
		if err != nil {
			return err
		}

		// index of expired_at field
		if err = c.EnsureIndex(mgo.Index{Key: []string{config.ExpiredAtField}, Name: config.ExpiredAtField}); err != nil {
			return err
		}

		// unique index of key field
		if err = c.EnsureIndex(mgo.Index{Key: []string{config.KeyField}, Name: config.KeyField, Unique: true}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return
	}

	return
}

func (p *MgoCache) Get(key string) *object.Value {
	values, ok := p.fetchValues(key)
	if !ok {
		return object.ValueTo(nil)
	}

	value, _ := p.getValue(values)
	return object.ValueTo(value)
}

func (p *MgoCache) Set(key string, val interface{}, params ...int) (err error) {
	p.connect(func(c *mgo.Collection) error {
		var timeout time.Duration
		if len(params) > 0 {
			switch {
			case params[0] > 0:
				timeout = time.Duration(params[0]) * time.Second
			case params[0] == 0:
				timeout = time.Hour * 24 * 9999
			}
		}

		if timeout < 0 {
			timeout = DefaultTimeout
		}

		expiredAt := time.Now().Add(timeout)

		_, err = c.Upsert(bson.M{
			p.config.KeyField: key,
		}, bson.M{
			"$set": bson.M{
				p.config.KeyField:       key,
				p.config.ValueField:     val,
				p.config.ExpiredAtField: expiredAt,
			},
		})
		return err
	})
	return
}

func (p *MgoCache) Delete(key string) (err error) {
	p.connect(func(c *mgo.Collection) error {
		err = c.Remove(bson.M{
			p.config.KeyField: key,
		})
		return err
	})
	return
}

func (p *MgoCache) Incr(key string, params ...int) (err error) {
	cnt := 1
	if len(params) > 0 {
		cnt = params[0]
	}

	_, ok := p.fetchValues(key)
	if !ok {
		return ErrMissKey
	}

	p.connect(func(c *mgo.Collection) error {
		err = c.Update(bson.M{
			p.config.KeyField: key,
		}, bson.M{
			"$inc": bson.M{
				p.config.ValueField: cnt,
			},
		})
		return err
	})
	return
}

func (p *MgoCache) Decr(key string, params ...int) (err error) {
	cnt := -1
	if len(params) > 0 {
		cnt = 0 - params[0]
	}

	_, ok := p.fetchValues(key)
	if !ok {
		return ErrMissKey
	}

	p.connect(func(c *mgo.Collection) error {
		err = c.Update(bson.M{
			p.config.KeyField: key,
		}, bson.M{
			"$inc": bson.M{
				p.config.ValueField: cnt,
			},
		})
		return err
	})
	return
}

func (p *MgoCache) Has(key string) bool {
	values, ok := p.fetchValues(key)
	if !ok {
		return false
	}

	_, ok = p.getValue(values)
	return ok
}

func (p *MgoCache) Clean() (err error) {
	p.connect(func(c *mgo.Collection) error {
		_, err = c.RemoveAll(nil)
		return err
	})
	return
}

func (p *MgoCache) GC() (err error) {
	p.connect(func(c *mgo.Collection) error {
		_, err = c.RemoveAll(bson.M{
			p.config.ExpiredAtField: bson.M{
				"$lte": time.Now(),
			},
		})
		return err
	})
	return
}

func (p *MgoCache) fetchValues(key string) (values map[string]interface{}, ok bool) {
	p.connect(func(c *mgo.Collection) error {
		c.Find(bson.M{
			p.config.KeyField: key,
		}).One(&values)
		return nil
	})
	return values, values != nil
}

func (p *MgoCache) getValue(values map[string]interface{}) (value interface{}, ok bool) {
	if values == nil {
		return
	}

	// get expired time
	expiredAt, exists := values[p.config.ExpiredAtField].(time.Time)
	if exists {
		// not expired yet
		if time.Now().Before(expiredAt) {

			// get cached value
			value, ok = values[p.config.ValueField]
		}
	}
	return
}
