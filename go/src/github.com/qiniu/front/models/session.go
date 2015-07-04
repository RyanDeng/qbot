package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	LoginUid     = "login_uid"
	LoginToken   = "login_token"
	LoginRefresh = "login_refresh"
	LoginExpired = "login_expired"
)

var (
	Session = &_Session{NewConnector(
		DefaultMgo,
		bson.M{
			"name": "session",
			"index": []string{
				"values." + LoginUid,
			},
			"unique": []string{
				"sid",
			},
		},
	)}
)

/*
type SessionModel struct {
	Id        bson.ObjectId          `bson:"_id"`
	Sid       string                 `bson:"sid"`
	Values    map[string]interface{} `bson:"values"`
	CreatedAt time.Time              `bson:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at"`
}
*/
type _Session struct{ Connector }

func (e *_Session) DestroyUserSession(uid uint32) (err error) {
	err = e.Invoke(func(c *mgo.Collection) error {
		_, er := c.RemoveAll(bson.M{
			"values." + LoginUid: uid,
		})
		return er
	})
	return
}
