package mail

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

const (
	mailgunAPITmpl = "https://api:%s@api.mailgun.net/v2/%s/messages"
)

type MailgunService struct {
	apiKey string

	from  string
	name  string
	reply string

	apiAddress string
}

func createMailgunService(apiKey, mailDomain, from, name, reply string) MailService {
	return &MailgunService{
		from:       from,
		name:       name,
		reply:      reply,
		apiAddress: fmt.Sprintf(mailgunAPITmpl, apiKey, mailDomain),
	}
}

func (service *MailgunService) buildForm(msg MailMessage) (url.Values, error) {
	params := url.Values{}

	params.Add("from", fmt.Sprintf("%s<%s>", msg.Name, msg.From))
	params.Add("h:Reply-To", msg.Reply)

	for _, v := range msg.To {
		params.Add("to", v)
	}

	for _, v := range msg.Bcc {
		params.Add("bcc", v)
	}

	for _, v := range msg.Cc {
		params.Add("cc", v)
	}

	params.Add("subject", msg.Subject)
	params.Add("html", msg.Content)
	//
	if msg.ExtraHeader != nil && len(msg.ExtraHeader) != 0 {
		for k, v := range msg.ExtraHeader {
			key := fmt.Sprintf("h:%s", k)
			params.Add(key, v)
		}
	}

	if msg.Options != nil {
		if mailList, ok := msg.Options["use_maillist"]; ok {
			params.Add("to", mailList)
		}
	}

	for _, v := range msg.Tag {
		params.Add("o:tag", v)
	}

	return params, nil
}

func checkMailgunRes(resp *http.Response) error {
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Read Mailgun res error %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mailgun error %d %s %s", resp.StatusCode, resp.Status, string(r))
	}

	return nil
}

func (service *MailgunService) Send(msg_ *MailMessage) (err error) {
	msg := *msg_

	if msg.Name == "" {
		msg.Name = service.name
	}
	if msg.From == "" {
		msg.From = service.from
	}
	if msg.Reply == "" {
		msg.Reply = service.reply
	}

	err = msg.Validate()
	if err != nil {
		return
	}

	params, err := service.buildForm(msg)
	if err != nil {
		return errors.New("build form fail " + err.Error())
	}

	client := http.Client{
		Transport: DefaultTransport,
	}

	resp, err := client.PostForm(service.apiAddress, params)
	if err != nil {
		return errors.New("PostForm fail " + err.Error())
	}

	return checkMailgunRes(resp)
}

func buildMailgunMultiPart(params url.Values, attachFiles []string) (contentType string, buf *bytes.Buffer, err error) {
	buf = bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)

	err = formToField(writer, params)
	if err != nil {
		return
	}
	for _, v := range attachFiles {
		err = attach(writer, "attachment", v)
		if err != nil {
			return
		}
	}

	err = writer.Close()
	if err != nil {
		return
	}
	return writer.FormDataContentType(), buf, nil
}

func (service *MailgunService) SendWithAttach(msg_ *MailMessage, attachFiles []string) (err error) {
	msg := *msg_

	if msg.Name == "" {
		msg.Name = service.name
	}
	if msg.From == "" {
		msg.From = service.from
	}
	if msg.Reply == "" {
		msg.Reply = service.reply
	}

	err = msg.Validate()
	if err != nil {
		return
	}

	params, err := service.buildForm(msg)
	if err != nil {
		return
	}
	contentType, buf, err := buildMailgunMultiPart(params, attachFiles)
	if err != nil {
		return
	}
	resp, err := http.Post(service.apiAddress, contentType, buf)
	if err != nil {
		return
	}
	return checkMailgunRes(resp)
}
