package sessions

import (
	"fmt"

	"hd.qiniu.com/utils/object"
)

const (
	SESSION_ID_LENGTH        = 32
	COOKIE_VALUE_SPLIT       = "," // value,value,value
	COOKIE_VALUE_PARTS_SPLIT = "|" // value1|value2|value3,name,time
)

var (
	ErrDuplicateSid    = fmt.Errorf("<Sessions> session id duplicated can not create")
	ErrNotFoundSession = fmt.Errorf("<Sessions> session not found")
	ErrEmptySecretKey  = fmt.Errorf("<Sessions> please set session secret key")

	_SESSION_ALPHABETS = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type SessionStore interface {
	Sid() string                       // back current sessionID
	Set(key string, value interface{}) // set session value
	Get(key string) *object.Value      // get session value
	Delete(key string)                 // delete session value
	Has(key string) bool               // check session key exist
	Values() map[string]interface{}    // duplicate all values
	Destroy() error                    // delete session in store
	Clean()                            // clean all data
	Flush() error                      // release the resource & save data to provider
}

type SessionProvider interface {
	Create(sid string, params ...map[string]interface{}) (SessionStore, error)
	Read(sid string) (SessionStore, error)
	Regenerate(oldsid, sid string) (SessionStore, error)
	Destroy(sid string) error
	GC() error // use for interval GC
	Config() *Config
}
