package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
