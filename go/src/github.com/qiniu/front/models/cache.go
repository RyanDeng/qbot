package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

var (
	Cache = &_Cache{NewConnector(
		DefaultMgo,
		bson.M{
			"name": "cache",
			"unique": []string{
				"key",
			},
		},
	)}
)

type CacheModel struct {
	Id        bson.ObjectId `bson:"_id"`
	Key       string        `bson:"key"`
	Value     interface{}   `bson:"value"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type _Cache struct{ Connector }
