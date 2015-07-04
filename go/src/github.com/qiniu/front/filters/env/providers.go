package env

import (
	"html/template"
	"path/filepath"

	"github.com/teapots/params"
	"github.com/teapots/render"
	"github.com/teapots/teapot"
	"hd.qiniu.com/env/global"
	"hd.qiniu.com/providers/account"
	"hd.qiniu.com/providers/biz"
	"hd.qiniu.com/providers/cache"
	"hd.qiniu.com/providers/mail"
	"hd.qiniu.com/providers/price"
	"hd.qiniu.com/providers/session"
	"hd.qiniu.com/utils/crypto"
	"hd.qiniu.com/utils/sessions"
)

func ConfigProviders(tea *teapot.Teapot) {
	templateDir := filepath.Join(tea.Config.RunPath, "./views")
	tea.Provide(render.Renderer(render.Options{
		Directory: templateDir,
		Delims: render.Delims{
			Left:  "[[",
			Right: "]]",
		},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"Sign": func(id, secret string) string {
					return crypto.MD5(id + secret)
				},
				"MD5": crypto.MD5,
				"Unescape": func(x string) interface{} {
					return template.HTML(x)
				},
				"Inc": func(n int) int {
					return n + 1
				},
				"Dec": func(n int) int {
					return n - 1
				},
				"Asset": func(subPath string) string {
					return global.Env.AssetPrefix + subPath
				},
			},
		},
	}))

	tea.Provide(params.ParamsParser())

	manager := session.SessionManager(sessions.Config{
		Logger: tea.Logger(),

		SecretKey:    global.Env.ClientSecret,
		CookieName:   global.Env.CookiePrefix + "_SESSION",
		CookieSecure: global.Env.CookieSecure,

		SessionExpire: 3600 * 24 * 7, // 7 days
		CookieExpire:  0,             // session life of browser

		AutoExpire: false,

		CookieRememberName: global.Env.CookiePrefix + "_REMEMBER",
		RememberExpire:     3600 * 24 * 7, // 7 days
	})

	global.SessionManager = manager
	tea.Provide(manager)
	tea.Provide(session.SessionStore())

	adminService := account.AdminService(
		global.Env.AdminAuth.AccountHost,
		global.Env.AdminAuth.ClientId,
		global.Env.AdminAuth.ClientSecret,
		global.Env.AdminAuth.UserName,
		global.Env.AdminAuth.Password,
	)
	tr := adminService.Transport()
	err := adminService.Auth()
	if err != nil {
		panic(err)
	}

	global.AdminService = adminService

	priceService := price.PriceService(global.Env.Price.Host, tr)
	tea.Provide(priceService)
	global.PriceService = priceService

	bizService := biz.BizService(global.Env.Biz.Host, tr)
	tea.Provide(bizService)
	global.BizService = bizService

	tea.Provide(account.AccountService(global.Env.PortalOAuth.AccountLink))
	tea.Provide(account.OAuthService(
		global.Env.PortalOAuth.AuthLink,
		global.Env.PortalOAuth.TokenLink,
		global.Env.PortalOAuth.ClientId,
		global.Env.PortalOAuth.ClientSecret,
	))

	tea.Provide(cache.MgoCache())

	mailService := mail.Mail(global.Env.Mail.ApiKey, global.Env.Mail.MailDomain, global.Env.Mail.From, global.Env.Mail.Name, global.Env.Mail.Reply, templateDir, tea.Logger())
	global.MailService = mailService
	tea.Provide(mailService)
}
