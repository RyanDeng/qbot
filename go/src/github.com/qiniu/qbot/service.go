package qbot

import (
	"fmt"
	"github.com/qiniu/postmans/interfaces"
	qmgo "github.com/qiniu/qbot/mgo"
	"strings"
)

type DBSettings struct {
	Mongodb *qmgo.Config `json:"mongodb"`

	ReminderColl string `json:"reminder_coll"`
	ContactColl  string `json:"contact_coll"`
}

type Config struct {
	DBSettings DBSettings
}

type Service struct {
	Config

	contactTbl  *ContactTbl
	reminderTbl *ReminderTbl
	Handles     []Handle
}

func NewService(cfg *Config) (*Service, error) {
	mcfg := cfg.DBSettings.Mongodb

	mcfg.Coll = cfg.DBSettings.ReminderColl
	reminderTbl, err := NewReminderTbl(mcfg)
	if err != nil {
		return nil, err
	}
	mcfg.Coll = cfg.DBSettings.ContactColl
	contactTbl, err := NewContactTbl(mcfg)
	if err != nil {
		return nil, err
	}

	handles := []Handle{
		ReminderHandle{},
		ContactHandle{},
	}

	p := &Service{
		Config: *cfg,

		reminderTbl: reminderTbl,
		contactTbl:  contactTbl,
		Handles:     handles,
	}

	return p, nil
}

func (s *Service) AI(msg *interfaces.Msg, pstman interfaces.Postman) {
	var matchHandle Handle
	var result string
	for _, handle := range s.Handles {
		if isMatchHandle(msg.Msg, handle.KeyWords()) {
			matchHandle = handle
			fmt.Println("match handle", handle.KeyWords())
		}
	}
	if matchHandle != nil {
		result = matchHandle.ThinkOut(s, msg.Msg)
	} else {
		result = "对不起我暂时有点笨，现在只能理解索要联系方式和设置提醒"
	}
	pstman.SendMsg(msg.From, result)
	return
}

func isMatchHandle(msg string, keyWords []string) bool {
	for _, keyWord := range keyWords {
		if strings.Contains(msg, keyWord) {
			return true
		}
	}
	return false
}

func (s *Service) GroupAI(grpmsg *interfaces.GroupMsg, pstman interfaces.Postman) {
	var matchHandle Handle
	var result string
	for _, handle := range s.Handles {
		if isMatchHandle(grpmsg.Msg, handle.KeyWords()) {
			matchHandle = handle
		}
	}
	if matchHandle != nil {
		result = matchHandle.GroupThinkOut(s, grpmsg.Msg)
	} else {
		result = "对不起我暂时有点笨，现在只能理解索要联系方式和设置提醒"
	}
	fmt.Println("group id ", grpmsg.GroupId)
	pstman.SendGroupMsg(grpmsg.GroupId, result)
	return
}
