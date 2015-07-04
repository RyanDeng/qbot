package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	_SendCloudSendAPI = "sendcloud.sohu.com/webapi/mail.send.json"
)

type sendCloudMailService struct {
	userName string
	password string

	from  string
	name  string
	reply string

	apiAddress string
}

func createSendCloud(user, password, from, name, reply, https string) MailService {
	var api string
	if https == "false" {
		api = "http://" + _SendCloudSendAPI
	} else {
		api = "https://" + _SendCloudSendAPI
	}
	return &sendCloudMailService{
		userName:   user,
		password:   password,
		from:       from,
		name:       name,
		reply:      reply,
		apiAddress: api,
	}
}

func (service *sendCloudMailService) buildForm(msg MailMessage) (url.Values, error) {
	params := url.Values{}

	params.Add("api_user", service.userName)
	params.Add("api_key", service.password)

	params.Add("from", msg.From)
	params.Add("fromname", msg.Name)
	params.Add("replyto", msg.Reply)

	params.Add("to", strings.Join(msg.To, ";"))
	if msg.Bcc != nil && len(msg.Bcc) > 0 {
		params.Add("bcc", strings.Join(msg.Bcc, ";"))
	}
	if msg.Cc != nil && len(msg.Cc) > 0 {
		params.Add("cc", strings.Join(msg.Cc, ";"))
	}

	params.Add("subject", msg.Subject)
	params.Add("html", msg.Content)

	if msg.ExtraHeader != nil && len(msg.ExtraHeader) != 0 {
		data, err := json.Marshal(msg.ExtraHeader)
		if err != nil {
			return nil, err
		}
		params.Add("headers", string(data))
	}

	if msg.Options != nil {
		if mailList, ok := msg.Options["use_maillist"]; ok {
			params.Add("use_maillist", mailList)
		}
	}

	if msg.Tag != nil && len(msg.Tag) != 0 {
		params.Add("label", msg.Tag[0])
	}
	return params, nil
}

type _SendCloudRet struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func checkSendCloudResp(resp *http.Response) error {
	r, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SendCloud error %d %s %s", resp.StatusCode, resp.Status, string(r))
	}
	ret := _SendCloudRet{}
	err := json.Unmarshal(r, &ret)
	if err != nil {
		return fmt.Errorf("SendCloud json response decode error %s", err.Error())
	}
	if ret.Message != "success" {
		return fmt.Errorf("SendCloud error %v", ret.Errors)
	}
	return nil
}

func (service *sendCloudMailService) Send(msg_ *MailMessage) (err error) {
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

	return checkSendCloudResp(resp)
}

func buildMultiPart(params url.Values, attachFiles []string) (contentType string, buf *bytes.Buffer, err error) {
	buf = bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)

	err = formToField(writer, params)
	if err != nil {
		return
	}
	for _, v := range attachFiles {
		err = attach(writer, "files", v)
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

func (service *sendCloudMailService) SendWithAttach(msg *MailMessage, attachFiles []string) (err error) {
	err = msg.Validate()
	if err != nil {
		return
	}

	params, err := service.buildForm(*msg)
	if err != nil {
		return
	}
	contentType, buf, err := buildMultiPart(params, attachFiles)
	if err != nil {
		return
	}
	resp, err := http.Post(service.apiAddress, contentType, buf)
	if err != nil {
		return
	}
	return checkSendCloudResp(resp)
}
