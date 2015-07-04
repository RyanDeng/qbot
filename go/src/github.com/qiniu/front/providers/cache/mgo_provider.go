package cache

import (
	"github.com/teapots/teapot"
	"hd.qiniu.com/models"
	"hd.qiniu.com/services/cache"
)

func MgoCache() interface{} {
	return func(log teapot.Logger) (mgoCache cache.Cache) {
		config := cache.NewMgoConfig()
		mgoCache, err := cache.NewMgoCache(config, models.Cache.Invoke)
		if err != nil {
			log.Error(err)
		}
		return
	}
}
