package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/teapots/params"
	"github.com/teapots/teapot"
	"hd.qiniu.com/env/global"
	"hd.qiniu.com/models"
	"hd.qiniu.com/services/account"
	"hd.qiniu.com/services/cache"
	"hd.qiniu.com/services/mail"
	"hd.qiniu.com/utils/crypto"
	"hd.qiniu.com/utils/object"
	"hd.qiniu.com/utils/sessions"
	"hd.qiniu.com/utils/verifycode"
)

const (
	flashMsgKey          = "flash_msg"
	loginedUidKey        = "LOGINED_UID"
	tokenCacheKey        = "OAUTH_TOKEN"
	tokenExpiresIn       = 86400 * 30
	verifyEmailCacheKey  = "VERIFY_EMAIL"
	verifyEmailExpiresIn = 86400
	verifyEmailPrefix    = "BPUy072zw2WGgFv65b4fZhjkTUdfcUF7"
	verifyCodeKey        = "VERIFY_CODE"
)

type VerifyEmail struct {
	Uid   uint32
	Code  string
	Email string
}

type Base struct {
	Req            *http.Request          `inject`
	Rw             http.ResponseWriter    `inject`
	Params         *params.Params         `inject`
	Log            teapot.Logger          `inject`
	Session        sessions.SessionStore  `inject`
	AccountService account.AccountService `inject`
	OAuthService   account.OAuthService   `inject`
	Cache          cache.Cache            `inject`
	MailService    mail.MailService       `inject`
}

func (c *Base) ExchangeOAuthToken(code string) (*account.OAuthToken, error) {
	return c.OAuthService.Exchange(code)
}

func (c *Base) GetOAuthToken(uid uint32) (token *account.OAuthToken, err error) {
	key := fmt.Sprintf("%s-%d", tokenCacheKey, uid)
	str := c.Cache.Get(key).String()
	if str == "" {
		err = account.ErrTokenNotFound
		return
	}

	token = &account.OAuthToken{}
	err = token.Unserilize(str)
	if err != nil {
		err = account.ErrTokenNotFound
		return
	}

	// 如果token过期，则刷新并持久化token
	if token.IsExpired() {
		err = account.ErrTokenExpired
		token, err = c.OAuthService.Refresh(token.RefreshToken)
		if err != nil {
			return
		}
		err = c.SetOAuthToken(uid, token)
	}

	if !token.IsValid() {
		err = account.ErrTokenInvalid
		return
	}

	return
}

func (c *Base) SetOAuthToken(uid uint32, token *account.OAuthToken) (err error) {
	str, err := token.Serilize()
	if err != nil {
		return
	}

	key := fmt.Sprintf("%s-%d", tokenCacheKey, uid)
	c.Cache.Set(key, str, tokenExpiresIn)
	return
}

func (c *Base) SetLoginUid(uid uint32) {
	c.Session.Set(loginedUidKey, uid)
}

func (c *Base) GetLoginUid() (uint32, error) {
	return c.Session.Get(loginedUidKey).Uint32()
}

func (c *Base) GetLoginUser() (*models.UserModel, error) {
	uid, err := c.GetLoginUid()
	if err != nil {
		return nil, err
	}
	if uid == 0 {
		return nil, errors.New("Not logined")
	}
	return models.User.FindByUid(uid)
}

func (c *Base) IsAdminLogin(activity string) bool {
	uids, ok := global.AdminUids[activity]
	if !ok {
		return false
	}
	uid, err := c.GetLoginUid()
	if err != nil {
		return false
	}
	return uids[uid]
}

func (c *Base) DestroySession() error {
	return c.Session.Destroy()
}

func (c *Base) AuthURL(state string) string {
	redirect := c.Scheme() + "://" + c.Req.Host + "/auth"
	return c.OAuthService.AuthURL(redirect, state)
}

func (c *Base) SendMail(tpl, to, subject string, val map[string]interface{}, tags ...string) error {
	content, err := mail.RenderMail(tpl, val)
	if err != nil {
		return err
	}

	return global.MailService.Send(&mail.MailMessage{
		To:      strings.Split(to, ","),
		Subject: subject,
		Content: content,
		Tag:     tags,
	})
}
func (c *Base) VerifyEmailURL(code string) string {
	return c.Scheme() + "://" + c.Req.Host + "/verify/email?code=" + code + "&sign=" + c.VerifyEmailSign(code)
}

func (c *Base) VerifyEmailSign(code string) string {
	return crypto.MD5(verifyEmailPrefix + code)
}

func (c *Base) SetVerifyEmailCode(code string, uid uint32) (err error) {
	key := fmt.Sprintf("%s-uid-%d", verifyEmailPrefix, uid)
	err = c.Cache.Set(key, code, verifyEmailExpiresIn)
	return
}

func (c *Base) GetVerifyEmailCode(uid uint32) (code string, err error) {
	key := fmt.Sprintf("%s-uid-%d", verifyEmailPrefix, uid)
	code = c.Cache.Get(key).String()
	if code == "" {
		err = cache.ErrMissKey
	}
	return
}

func (c *Base) SetVerifyEmail(code string, uid uint32, email string) (err error) {
	key := fmt.Sprintf("%s-%s", verifyEmailCacheKey, code)
	str, err := object.Serilize(&VerifyEmail{
		Uid:   uid,
		Code:  code,
		Email: email,
	})
	if err != nil {
		return
	}
	err = c.Cache.Set(key, str, verifyEmailExpiresIn)
	return
}

func (c *Base) GetVerifyEmail(code string) (data *VerifyEmail, err error) {
	key := fmt.Sprintf("%s-%s", verifyEmailCacheKey, code)
	str := c.Cache.Get(key).String()
	if str == "" {
		err = cache.ErrMissKey
		return
	}

	data = &VerifyEmail{}
	err = object.Unserilize(str, data)
	return
}

func (c *Base) RemoveVerifyEmail(code string) (err error) {
	key := fmt.Sprintf("%s-%s", verifyEmailCacheKey, code)
	err = c.Cache.Delete(key)
	return
}

// Gnerate a new image captcha code
func (c *Base) RenderVerifyCode(id string, width, height int, resp http.ResponseWriter) (err error) {
	s, err := verifycode.NewLen(4)
	if err != nil {
		return
	}
	code := ""
	data := []byte(s)
	for v := range data {
		data[v] %= 10
		code += strconv.FormatInt(int64(data[v]), 32)
	}

	key := fmt.Sprintf("%s-%s", verifyCodeKey, id)
	c.Session.Set(key, code)
	verifycode.NewImage(data, width, height).WriteTo(resp)

	return
}

func (c *Base) IsVerifyCodeMatch(id string, code string) bool {
	key := fmt.Sprintf("%s-%s", verifyCodeKey, id)
	verifyCode := c.Session.Get(key).String()
	c.Session.Delete(key)
	return (code == verifyCode)
}

func (c *Base) FlashMsg(msg string) {
	c.Session.Set(flashMsgKey, msg)
}

func (c *Base) GetFlashMsg() (msg string) {
	msg = c.Session.Get(flashMsgKey).String()
	c.Session.Delete(flashMsgKey)
	return
}

func (c *Base) GetIP() (ip string) {
	ip = c.Req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Req.Header.Get("X-Real-IP")
	}
	return
}

func (c *Base) Scheme() (scheme string) {
	scheme = c.Req.Header.Get("X-Scheme")
	if scheme == "" {
		scheme = "http"
	}
	return
}

func (c *Base) GetDailyLimitTimes(key string) int {
	key = fmt.Sprintf("LIMIT_TIMES-%s", key)
	return c.Cache.Get(key).MustInt()
}

func (c *Base) SetDailyLimitTimes(key string, times int) (err error) {
	var (
		now              = time.Now()
		year, month, day = now.Date()
		expireTime       = time.Date(year, month, day+1, 0, 0, 0, 0, time.UTC)
		expiresIn        = int(expireTime.Unix() - now.Unix())
	)
	key = fmt.Sprintf("LIMIT_TIMES-%s", key)
	err = c.Cache.Set(key, times, expiresIn)
	return
}
