package sessions

import (
	"time"

	"github.com/teapots/teapot"
)

type Config struct {
	Logger teapot.Logger // logger

	SecretKey string // secure secret key

	CookieName   string // session cookie name
	CookieSecure bool   // is cookie use https?

	SessionExpire int // session expire seconds
	CookieExpire  int // session cookie expire seconds

	AutoExpire bool // is provider support auto expire?

	CookieRememberName string // hashed value of user for auto login

	RememberExpire int // auto login remember expire seconds
}

func (c *Config) SessionExpireSeconds() time.Duration {
	return time.Duration(c.SessionExpire) * time.Second
}

func (c *Config) RememberExpireSeconds() time.Duration {
	return time.Duration(c.RememberExpire) * time.Second
}
