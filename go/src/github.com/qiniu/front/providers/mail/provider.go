package mail

import (
	"github.com/teapots/teapot"
	"hd.qiniu.com/services/mail"
)

func Mail(apiKey, mailDomain, from, name, reply string, templateDir string, logger teapot.Logger) mail.MailService {
	mail.TemplateDir = templateDir
	config := mail.MailConfig{
		"Type":       "mailgun",
		"ApiKey":     apiKey,
		"MailDomain": mailDomain,
		"From":       from,
		"Name":       name,
		"Reply":      reply,
	}
	service, err := mail.NewMailService(config)
	if err != nil {
		logger.Error(err)
	}
	return service
}
