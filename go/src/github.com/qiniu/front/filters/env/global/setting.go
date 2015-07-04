package global

import (
	"github.com/teapots/teapot"
	"hd.qiniu.com/services/account"
	"hd.qiniu.com/services/biz"
	"hd.qiniu.com/services/mail"
	"hd.qiniu.com/services/price"
	"hd.qiniu.com/utils/sessions"
)

var (
	Env            *Setting
	SessionManager *sessions.SessionManager
	PriceService   price.PriceService
	BizService     biz.BizService
	MailService    mail.MailService
	AdminService   account.AdminService

	AdminUids = map[string]map[uint32]bool{}
)

type Setting struct {
	Teapot       *teapot.Teapot `conf:"-"`
	Config       *teapot.Config `conf:"-"`
	CookiePrefix string         `conf:"cookie_prefix"`
	CookieSecure bool           `conf:"cookie_secure"`
	ClientSecret string         `conf:"client_secret"`
	AssetPrefix  string         `conf:"asset_prefix"`
	Price        struct {
		Host string `conf:"host"`
	} `conf:"price"`
	Biz struct {
		Host string `conf:"host"`
	} `conf:"biz"`
	AdminAuth struct {
		AccountHost  string `conf:"account_host"`
		ClientId     string `conf:"admin_client_id"`
		ClientSecret string `conf:"admin_client_secret"`
		UserName     string `conf:"admin_username"`
		Password     string `conf:"admin_password"`
	} `conf:"admin_auth"`
	Mail struct {
		ApiKey     string `conf:"api_key"`
		MailDomain string `conf:"mail_domain"`
		From       string `conf:"from"`
		Name       string `conf:"name"`
		Reply      string `conf:"reply"`
	} `conf:"mail"`

	Mongo struct {
		Default string `conf:"db.default"`
	} `conf:"mongo"`

	PortalOAuth struct {
		ClientId     string `conf:"client_id"`
		ClientSecret string `conf:"client_secret"`
		AuthLink     string `conf:"auth_link"`
		TokenLink    string `conf:"token_link"`
		AccountLink  string `conf:"account_link"`
	} `conf:"portal_oauth"`

	Limit struct {
		MaxVerifyTimes         int `conf:"max_verify_times"`
		MaxDailySendEmailTimes int `conf:"max_daily_send_email_times"`
	} `conf:"limit"`
	Job struct {
		Switch string `conf:"switch"`
	} `conf:"job"`
}
