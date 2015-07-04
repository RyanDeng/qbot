package qbot

import (
	"hd.qiniu.com/models"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	Qbot = &_Qbot{
		Connector: models.NewConnector(
			models.DefaultMgo,
			bson.M{
				"name": "contact",
				"index": []string{
					"status",
				},
				"unique": []string{
					"phone",
				},
			},
		),
	}
)

type Status int

const (
	StatusWait Status = iota
	StatusSuccess
	StatusInjected
)

func (s Status) Humanize() string {
	switch s {
	case StatusSuccess:
		return "已通过"
	case StatusInjected:
		return "已拒绝"
	case StatusWait:
		return "待审核"
	}
	return ""
}

func (s Status) Int() int {
	return int(s)
}

type QbotModel struct {
	Id          bson.ObjectId `json:"_id"`
	Name        string        `json:"" bson:"name"`
	Phone       int64         `json:"" bson:"phone"`
	NickName    string        `json:"" bson:"nickname"`
	Email       string        `json:"email" bson:"email"`
	QQ          string        `json:"qq" bson:"qq"`
	Department  string        `json:"department" bson:"department"`
	CreatedTime time.Time     `json:"created_time"`
	UpdatedTime time.Time     `json:"updated_time"`
}

type _Qbot struct {
	models.Connector
}

func (p *_Qbot) New(name, nickname, email, qq, department string, phone int64) *QbotModel {
	return &QbotModel{
		Id:          bson.NewObjectId(),
		Name:        name,
		Phone:       phone,
		NickName:    nickname,
		Email:       email,
		QQ:          qq,
		Department:  department,
		CreatedTime: bson.Now(),
	}
}
