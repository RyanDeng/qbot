package qbot

import (
	"labix.org/v2/mgo/bson"
)

//-----------------
type From struct {
	CliType string `json:"cli_type" bson:"cli_type"`
	AccId   string `json:"account" bson:"account"`
}
type To struct {
	CliType string   `json:"cli_type" bson:"cli_type"`
	AccIds  []string `json:"accounts" bson:"accounts"`
}

type Reminder struct {
	Id        bson.ObjectId `json:"-" bson:"_id"`
	From      From          `json:"from" bson:"from"`
	Tos       []To          `json:"to" bson:"to"`
	Time      int64         `json:"time" bson:"time"`
	Status    string        `json:"status" bson:"status"`
	Event     string        `json:"event" bson:"event"`
	CreatedAt int64         `json:"created_at bson:"created_at"`
}

//------------------
type Contact struct {
	Name       string   `json:"name" bson:"name"`
	Photo      string   `json:"photo" bson:"photo"`
	Phone      int64    `json:"phone" bson:"phone"`
	NickName   []string `json:"nickname" bson:"nickname"`
	Email      string   `json:"email" bson:"email"`
	QQ         string   `json:"qq" bson:"qq"`
	Department string   `json:"dept" bson:"dept"`
	CreatedAt  int64    `json:"created_at" bson:"created_at"`
	UpdatedAt  int64    `json:"updated_at" bson:"updated_at"`
}
