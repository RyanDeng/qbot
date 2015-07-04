package rand

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"
)

var (
	intChars = []byte("0123456789")
)

func RandomIntString(n int) (str string, err error) {
	b, err := randomCreateBytes(n, intChars...)
	if err != nil {
		return
	}
	str = string(b)
	return
}

func RandomString(n int) (str string, err error) {
	b, err := randomCreateBytes(n)
	if err != nil {
		return
	}
	str = string(b)
	return
}

func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	src := mrand.NewSource(time.Now().UnixNano())
	return mrand.New(src).Intn(max-min) + min
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
