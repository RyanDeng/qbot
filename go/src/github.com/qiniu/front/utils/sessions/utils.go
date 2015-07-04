package sessions

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func CreateSid() (sid string) {
	var (
		retried int
	)

retry:
	bytes, err := randomCreateBytes(SESSION_ID_LENGTH, _SESSION_ALPHABETS...)
	if err != nil {
		retried += 1
		if retried > 3 {
			// 随机字符串不能创建，Critical Error
			panic(err)
		}
		time.Sleep(time.Second * 1)
		goto retry
	}

	sid = string(bytes)
	return
}

func DecodeSecureValue(value string, secretKey string) (raw string, createdAt time.Time, ok bool) {
	rawBytes, _ := base64.URLEncoding.DecodeString(value)
	value = string(rawBytes)

	parts := strings.SplitN(value, COOKIE_VALUE_SPLIT, 3)
	if len(parts) < 3 {
		return
	}

	vRaw := strings.TrimSpace(parts[0])
	vCreated := strings.TrimSpace(parts[1])
	vHash := strings.TrimSpace(parts[2])

	if vRaw == "" || vCreated == "" || vHash == "" {
		return
	}

	vTime, _ := strconv.ParseInt(vCreated, 10, 64)
	if vTime <= 0 {
		return
	}

	vRaw, err := url.QueryUnescape(vRaw)
	if err != nil {
		return
	}

	h := hmac.New(sha1.New, []byte(secretKey))
	_, err = h.Write([]byte(vRaw + vCreated))
	if err != nil {
		return
	}

	if hex.EncodeToString(h.Sum(nil)) != vHash {
		return
	}

	raw = vRaw
	createdAt = time.Unix(0, vTime)
	ok = true
	return
}

func EncodeSecureValue(raw string, secretKey string, createdAt time.Time) (value string, ok bool) {
	if raw == "" || secretKey == "" {
		return
	}

	timeValue := strconv.FormatInt(createdAt.UnixNano(), 10)

	h := hmac.New(sha1.New, []byte(secretKey))
	if _, err := h.Write([]byte(raw + timeValue)); err != nil {
		return
	}

	hash := hex.EncodeToString(h.Sum(nil))

	src := url.QueryEscape(raw)
	src += COOKIE_VALUE_SPLIT + timeValue + COOKIE_VALUE_SPLIT + hash

	value = base64.URLEncoding.EncodeToString([]byte(src))
	ok = true
	return
}

// randomCreateBytes generate random []byte by specify chars.
func randomCreateBytes(n int, alphabets ...byte) ([]byte, error) {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	if num, err := rand.Read(bytes); num != n || err != nil {
		if err == nil {
			err = fmt.Errorf("random string not enough length: need %d but %d", n, num)
		}
		return nil, err
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return bytes, nil
}

func isExpired(createdAt time.Time, expired time.Duration) bool {
	return createdAt.Add(expired).Before(time.Now())
}

// get last same name cookie from cookies
// http://play.golang.org/p/LDfjMnJnhI
func getCookie(cookies []*http.Cookie, name string) (cookie *http.Cookie, ok bool) {
	for i := len(cookies) - 1; i >= 0; i-- {
		if cookies[i].Name == name {
			ok = true
			cookie = cookies[i]
			break
		}
	}
	return
}
