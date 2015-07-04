package cache

import (
	"errors"
	"time"

	"hd.qiniu.com/utils/object"
)

const (
	DefaultTimeout = time.Minute
)

var (
	ErrMissKey = errors.New("err_missed_key")
)

type Cache interface {
	Get(key string) *object.Value                         // get cached value by key
	Set(key string, val interface{}, params ...int) error // set cached key, value with optional timeout seconds
	Delete(key string) error                              // delete cached value by key
	Incr(key string, params ...int) error                 // incr integer value
	Decr(key string, params ...int) error                 // decr integer value
	Has(key string) bool                                  // check cached key exists
	Clean() error                                         // clean all cached values
	GC() error                                            // use for interval GC
}
