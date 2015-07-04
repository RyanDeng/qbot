package sessions

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeSecureValue(t *testing.T) {
	_, ok := EncodeSecureValue("sid", "", time.Now())
	assert.False(t, ok)

	_, ok = EncodeSecureValue("", "secretKey", time.Now())
	assert.False(t, ok)

	value, ok := EncodeSecureValue("sid", "secretKey", time.Now())
	assert.True(t, ok)

	assert.NotEmpty(t, value)
}

func Test_DecodeSecureValue(t *testing.T) {
	created := time.Now()
	value, ok := EncodeSecureValue("sid", "secretKey", created)
	assert.True(t, ok)
	assert.NotEmpty(t, value)

	_, vTime, ok := DecodeSecureValue("", "secretKey")
	assert.False(t, ok)
	assert.True(t, vTime.IsZero())

	sid, vTime, ok := DecodeSecureValue("sid", "")
	assert.False(t, ok)
	assert.True(t, vTime.IsZero())

	sid, vTime, ok = DecodeSecureValue(value, "secretKey")
	assert.True(t, ok)
	assert.Equal(t, "sid", sid)
	assert.Equal(t, vTime.UnixNano(), created.UnixNano())

	rawBytes, _ := base64.StdEncoding.DecodeString(value)

	assert.Contains(t, string(rawBytes), COOKIE_VALUE_SPLIT)

	assert.Equal(t, len(strings.Split(string(rawBytes), COOKIE_VALUE_SPLIT)), 3)

	sid, vTime, ok = DecodeSecureValue(value, "wrong secret key")
	assert.False(t, ok)
	assert.Empty(t, sid)
	assert.True(t, vTime.IsZero())

	sid, vTime, ok = DecodeSecureValue("wrong"+COOKIE_VALUE_SPLIT+"wrong", "secretKey")
	assert.False(t, ok)
	assert.Empty(t, sid)
	assert.True(t, vTime.IsZero())
}

func Test_SecureValueEscapeChars(t *testing.T) {
	raw := "A|B,C"

	created := time.Now()
	value, ok := EncodeSecureValue(raw, "secretKey", created)
	assert.True(t, ok)
	assert.NotEmpty(t, value)

	vRaw, _, ok := DecodeSecureValue(value, "secretKey")
	assert.True(t, ok)
	assert.Equal(t, raw, vRaw)
}
