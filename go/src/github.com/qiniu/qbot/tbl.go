package qbot

import (
	"labix.org/v2/mgo"
	"time"
	//"labix.org/v2/mgo/bson"
	qmgo "github.com/qiniu/qbot/mgo"
)

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
