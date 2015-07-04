package sessions

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SessionManager struct {
	provider SessionProvider
}

func NewSessionManager(provider SessionProvider) *SessionManager {
	if provider.Config().SecretKey == "" {
		// SecretKey 为空可能导致安全问题，坚决 panic
		panic(ErrEmptySecretKey)
	}

	manager := new(SessionManager)
	manager.provider = provider
	return manager
}

func (m *SessionManager) GC(intervals ...time.Duration) {
	var interval time.Duration
	if len(intervals) > 0 {
		interval = intervals[0]

		// 小于1分钟的，默认回收时间设为1小时
		if interval < time.Minute {
			interval = time.Hour
		}
	}

	// 启动 Provider 的 GC
	err := m.provider.GC()
	if err != nil {
		log.Println("[SessionManager.GC] err:", err)
	}

	time.AfterFunc(interval, func() {
		// 时间到，再次执行 GC
		m.GC(interval)
	})
}

// 开始一个Session，从请求中获取Sid，或者创建一个新的
func (m *SessionManager) Start(w http.ResponseWriter, r *http.Request) (sess SessionStore, createdAt time.Time, err error) {
	sid, createdAt, ok := m.readSidFromRequest(r)
	if !ok {
		sess, err = m.createSession()
		if err != nil {
			return
		}

		createdAt = m.WriteSessionCookie(r, w, sess.Sid())

	} else {
		sess, err = m.provider.Read(sid)
		if err == ErrNotFoundSession {
			sess, err = m.createSession()
			if err != nil {
				return
			}

			createdAt = m.WriteSessionCookie(r, w, sess.Sid())
		}
	}
	return
}

// 删除当前的session
func (m *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) error {
	var (
		config = m.provider.Config()
	)

	sid, _, ok := m.readSidFromRequest(r)
	if !ok {
		return nil
	}

	// error can secure skip
	_ = m.provider.Destroy(sid)

	cookie := &http.Cookie{
		Name:     config.CookieName,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
		Expires:  time.Now(),
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)

	return nil
}

// 为现有的session数据更换sid
func (m *SessionManager) Regenerate(w http.ResponseWriter, r *http.Request, params ...map[string]interface{}) (sess SessionStore, err error) {
	// 获取当前的sid，未找到时创建新的
	oldsid, _, ok := m.readSidFromRequest(r)
	if !ok {
		sess, err = m.createSession(params...)
		if err != nil {
			return
		}

		m.WriteSessionCookie(r, w, sess.Sid())
		return
	}

	var (
		sid     string
		retried int
	)

retry:
	sid = CreateSid()
	sess, err = m.provider.Regenerate(oldsid, sid)

	switch err {
	case ErrNotFoundSession:
		// 未找到时创建新的
		sess, err = m.createSession(params...)
		if err != nil {
			return
		}

	case ErrDuplicateSid:
		// 遇到重复的再次重试
		retried += 1
		if retried >= 3 {
			return
		}
		goto retry
	}

	m.WriteSessionCookie(r, w, sess.Sid())
	return
}

// 设置 Session Cookie
func (m *SessionManager) WriteSessionCookie(r *http.Request, w http.ResponseWriter, sid string) (createdAt time.Time) {
	var (
		config = m.provider.Config()
	)

	createdAt = time.Now()

	// secure cookie value of sid
	value, ok := EncodeSecureValue(sid, config.SecretKey, createdAt)
	if !ok {
		return
	}

	cookie := &http.Cookie{
		Name:     config.CookieName,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.CookieSecure,
	}

	if config.CookieExpire >= 0 {
		cookie.MaxAge = config.CookieExpire
	}

	http.SetCookie(w, cookie)
	r.AddCookie(cookie)

	return
}

// 创建新 Session
func (m *SessionManager) createSession(params ...map[string]interface{}) (sess SessionStore, err error) {
	var (
		sid     string
		retried int
	)

retry:
	sid = CreateSid()
	sess, err = m.provider.Create(sid, params...)
	if err == ErrDuplicateSid {
		retried += 1
		if retried >= 3 {
			return
		}
		goto retry
	}
	return
}

func (m *SessionManager) validSid(value string) (sid string, ok bool) {
	value = strings.TrimSpace(value)

	if len(value) != SESSION_ID_LENGTH {
		return
	}

	return value, true
}

func (m *SessionManager) readSidFromRequest(r *http.Request) (sid string, created time.Time, ok bool) {
	var (
		config = m.provider.Config()
	)

	cookie, ok := getCookie(r.Cookies(), config.CookieName)
	if !ok {
		return
	}

	value, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return
	}

	value, vTime, ok := DecodeSecureValue(value, config.SecretKey)
	if !ok {
		return
	}

	value, ok = m.validSid(value)
	if !ok {
		return
	}

	sid = value
	created = vTime
	ok = true
	return
}
