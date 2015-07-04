package qbot

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
	//"labix.org/v2/mgo/bson"
	qmgo "github.com/qiniu/qbot/mgo"
)

type M bson.M

type ContactTbl struct {
	coll *mgo.Collection
}

func NewContactTbl(cfg *qmgo.Config) (*ContactTbl, error) {
	s := qmgo.Open(cfg)

	return &ContactTbl{
		coll: s.Coll,
	}, nil
}

func (t *ContactTbl) Insert(contact *Contact) error {
	c := qmgo.CopyCollection(t.coll)
	defer qmgo.CloseCollection(c)

	contact.CreatedAt = time.Now().UnixNano()

	return c.Insert(contact)
}
func (t *ContactTbl) SearchByAllName(name string) (contacts []Contact, err error) {
	c := qmgo.CopyCollection(t.coll)
	defer qmgo.CloseCollection(c)

	sel := M{
		"$or": []M{
			M{"name": name},
			M{"nickname": name},
		},
	}
	fmt.Println("search condition", sel)

	err = c.Find(sel).All(&contacts)
	fmt.Println(contacts)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return
}

//-----------------
type ReminderTbl struct {
	coll *mgo.Collection
}

func NewReminderTbl(cfg *qmgo.Config) (*ReminderTbl, error) {
	s := qmgo.Open(cfg)

	return &ReminderTbl{
		coll: s.Coll,
	}, nil
}

func (t *ReminderTbl) Insert(reminder *Reminder) error {
	c := qmgo.CopyCollection(t.coll)
	defer qmgo.CloseCollection(c)

	reminder.CreatedAt = time.Now().UnixNano()

	return c.Insert(reminder)
}
