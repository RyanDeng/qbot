package models

import (
	"time"

	"hd.qiniu.com/models/enums"
	"labix.org/v2/mgo/bson"
)

var (
	User = &_User{
		Connector: NewConnector(
			DefaultMgo,
			bson.M{
				"name": "user",
				"index": []string{
					"uid,status",
					"email",
				},
				"unique": []string{
					"uid",
				},
			},
		),
	}
)

type UserModel struct {
	Id         bson.ObjectId `bson:"_id"`
	Uid        uint32        `bson:"uid"`
	Email      string        `bson:"email"`
	FullName   string        `bson:"fullname"`
	Gender     enums.Gender  `bson:"gender"`
	Status     enums.Status  `bson:"status"`
	CreateTime time.Time     `bson:"create_time"`
}

type _User struct {
	Connector
}

func (d *_User) New(uid uint32, email, fullName string, gender int) *UserModel {
	return &UserModel{
		Id:         bson.NewObjectId(),
		Uid:        uid,
		Email:      email,
		FullName:   fullName,
		Gender:     enums.Gender(gender),
		Status:     enums.StatusEnabled,
		CreateTime: time.Now(),
	}
}

func (d *_User) FindByUid(uid uint32) (user *UserModel, err error) {
	user = &UserModel{}
	err = d.Find(bson.M{
		"uid": uid,
	}, user)
	return
}

func (d *_User) FindByEmail(email string) (user *UserModel, err error) {
	user = &UserModel{}
	err = d.Find(bson.M{
		"email": email,
	}, user)
	return
}
