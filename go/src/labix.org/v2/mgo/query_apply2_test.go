package mgo

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"labix.org/v2/mgo/bson"
)

func TestQuery_Apply2(t *testing.T) {
	session, err := Dial("localhost")
	assert.NoError(t, err)

	col := session.DB("qbox_test").C("counters")
	col.RemoveId("user_userid")
	err = col.Insert(bson.M{
		"_id": "user_userid",
		"n":   int64(0),
	})
	assert.NoError(t, err)
	change := Change{
		Update:    bson.M{"$inc": bson.M{"n": 1}},
		ReturnNew: true,
	}
	res := map[string]interface{}{}

	_, err = col.FindId("user_userid").Apply2(change, &res)
	assert.NoError(t, err)
}
