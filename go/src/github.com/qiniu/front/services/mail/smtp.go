package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

import (
	message "github.com/sloonz/go-mime-message"
)

type smtpService struct {
	Port    int    `json:"Port"`
	Address string `json:"Address"`

	Account  string `json:"Account"`
	Password string `json:"Password"`

	From  string
	Name  string
	Reply string
}

func (m *smtpService) Send(msg_ *MailMessage) error {
	msg := *msg_
	if msg.From == "" {
		msg.From = m.From
	}
	if msg.Name == "" {
		msg.Name = m.Name
	}
	if msg.Reply == "" {
		msg.Reply = m.Reply
	}

	err := msg.Validate()
	if err != nil {
		return err
	}
	a := smtp.PlainAuth("", m.Account, m.Password, "")
	var crlf = "\r\n"
	var buf bytes.Buffer
	allTo := msg.To
	if msg.Reply != "" {
		fmt.Fprintf(&buf, "Reply-To: %s%s", msg.Reply, crlf)
	}
	fmt.Fprintf(&buf, "Date: %s%s", time.Now().Format(time.RFC822Z), crlf)
	fmt.Fprintf(&buf, "From: %s <%s>%s", message.EncodeWord(msg.Name), msg.From, crlf)
	fmt.Fprintf(&buf, "To: %s%s", strings.Join(msg.To, ","), crlf)
	if len(msg.Cc) > 0 {
		fmt.Fprintf(&buf, "Cc: %s%s", msg.Cc, crlf)
		allTo = append(allTo, msg.Cc...)
	}

	if len(msg.Bcc) > 0 {
		//fmt.Fprintf(&buf, "Bcc: %s%s", bcc, crlf)
		allTo = append(allTo, msg.Bcc...)
	}

	//Content-Type: multipart/alternative; boundary=047d7b33d6d615d96d04cc67487b
	fmt.Fprintf(&buf, "Subject: %s%s", message.EncodeWord(msg.Subject), crlf)
	fmt.Fprint(&buf, "Content-Type: multipart/alternative; boundary=047d7b33d6d615d96d04cc67487b", crlf, crlf)
	fmt.Fprint(&buf, "--047d7b33d6d615d96d04cc67487b", crlf)
	fmt.Fprint(&buf, "Content-Type: text/html; charset=UTF-8", crlf)

	fmt.Fprint(&buf, "Content-Transfer-Encoding: base64", crlf, crlf)
	body := base64.StdEncoding.EncodeToString([]byte(msg.Content))
	for idx := 0; ; idx += 76 {
		if idx+76 >= len(body) {
			fmt.Fprint(&buf, body[idx:], crlf)
			break
		} else {
			fmt.Fprint(&buf, body[idx:idx+76], crlf)
		}
	}
	fmt.Fprint(&buf, "--047d7b33d6d615d96d04cc67487b--", crlf)

	return m.send(m.Address+":"+strconv.Itoa(m.Port), a, msg.From, allTo, buf.Bytes())
}

func (m *smtpService) send(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}

	if ok, _ := c.Extension("AUTH"); a != nil && ok {
		resp := []byte("\x00" + m.Account + "\x00" + m.Password)
		str := base64.StdEncoding.EncodeToString(resp)
		_, _, err := m.cmd(c, 0, "AUTH %s %s", "PLAIN", str)
		if err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func (m *smtpService) cmd(c *smtp.Client, expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.Text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
	code, msg, err := c.Text.ReadResponse(expectCode)
	return code, msg, err
}

// todo, use multipart
func (service *smtpService) SendWithAttach(msg *MailMessage, attachFile []string) error {
	panic("no implementation")
	return nil
}
