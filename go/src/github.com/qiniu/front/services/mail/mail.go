package mail

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var DefaultTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
	Dial: (&net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}).Dial,
}

type MailService interface {
	Send(msg *MailMessage) error
	SendWithAttach(msg *MailMessage, attachFile []string) error
}

func NewMailService(mailConfig MailConfig) (res MailService, err error) {
	type_, ok1 := mailConfig["Type"]
	if !ok1 {
		err = errors.New("no type specified")
	}

	if type_ == "sendcloud" {
		user, ok1 := mailConfig["User"]
		password, ok2 := mailConfig["Password"]
		from := mailConfig["From"]
		name := mailConfig["Name"]
		reply := mailConfig["Reply"]
		https := mailConfig["https"]
		if !ok1 || !ok2 {
			err = errors.New("User and Password required")
			return
		}
		res = createSendCloud(user, password, from, name, reply, https)
		return
	}

	if type_ == "mailgun" {
		apiKey, ok1 := mailConfig["ApiKey"]
		mailDomain, ok2 := mailConfig["MailDomain"]
		if !ok1 || !ok2 {
			err = errors.New("Apikey and MailDomain is required")
			return
		}

		from := mailConfig["From"]
		name := mailConfig["Name"]
		reply := mailConfig["Reply"]

		res = createMailgunService(apiKey, mailDomain, from, name, reply)
		return
	}

	if type_ == "smtp" {
		host, ok1 := mailConfig["Host"]
		if !ok1 {
			err = errors.New("Host required")
			return
		}
		portStr, _ := mailConfig["Port"]
		if portStr == "" {
			portStr = "25"
		}
		var port int64
		port, err = strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			return
		}
		user, _ := mailConfig["User"]
		password, _ := mailConfig["Password"]
		from, _ := mailConfig["From"]
		name, _ := mailConfig["Name"]
		reply, _ := mailConfig["Reply"]

		res = &smtpService{
			Address:  host,
			Port:     int(port),
			Account:  user,
			Password: password,
			From:     from,
			Name:     name,
			Reply:    reply,
		}
		return
	}

	err = errors.New("unknown type: " + type_)
	return
}

type MailMessage struct {
	Name  string
	From  string
	Reply string

	To  []string
	Cc  []string
	Bcc []string

	Tag         []string
	ExtraHeader map[string]string

	Subject string
	Content string

	Options map[string]string
}

func (msg *MailMessage) Validate() error {
	// 内部使用，不检查from和to的邮件地址有效性
	if msg.Name == "" {
		return errors.New("no name")
	}

	if msg.From == "" {
		return errors.New("no from")
	}

	if msg.To == nil || len(msg.To) == 0 {
		return errors.New("no to")
	}

	if msg.Subject == "" {
		return errors.New("no subject")
	}

	if msg.Content == "" {
		return errors.New("no content")
	}
	return nil
}

func formToField(writer *multipart.Writer, params url.Values) (err error) {
	for k, v := range params {
		if len(v) == 0 {
			continue
		}
		err = writer.WriteField(k, v[0])
		if err != nil {
			return err
		}
	}
	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func CreateFormFile(w *multipart.Writer, fieldname, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return w.CreatePart(h)
}

func attach(writer *multipart.Writer, fieldname, file string) error {
	filename := path.Base(file)
	mimeType := mime.TypeByExtension(path.Ext(file))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w, err := CreateFormFile(writer, fieldname, filename, mimeType)
	if err != nil {
		return err
	}
	reader, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, reader)
	return err
}

type MailConfig map[string]string

const _Retry = 3

// please use utils.RevelOptionsMap to get config
func SendMail(config map[string]string, message MailMessage) (err error) {
	service, err := NewMailService(config)
	if err != nil {
		return
	}

	for i := 0; i < _Retry; i++ {
		err = service.Send(&message)
		if err == nil {
			break
		}
	}

	return err
}
