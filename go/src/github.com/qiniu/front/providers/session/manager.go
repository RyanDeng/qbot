package session

import (
	"hd.qiniu.com/models"
	"hd.qiniu.com/utils/sessions"
)

func SessionManager(config sessions.Config) *sessions.SessionManager {
	provider := sessions.NewMgoProvider(config, models.Session.Invoke)
	manager := sessions.NewSessionManager(provider)
	return manager
}
