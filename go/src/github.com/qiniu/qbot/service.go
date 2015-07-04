package qbot

import (
	"github.com/qiniu/postmans/interfaces"
	qmgo "github.com/qiniu/qbot/mgo"
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

	p := &Service{
		Config: *cfg,

		reminderTbl: reminderTbl,
		contactTbl:  contactTbl,
	}

	return p, nil
}

func (s *Service) AI(msg interfaces.Msg) {

}

func (s *Service) GroupAI(grpmsg interfaces.GroupMsg) {

}
