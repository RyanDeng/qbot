package env

import (
	"time"

	"github.com/teapots/teapot"
	"hd.qiniu.com/env/global"
	"hd.qiniu.com/models"
	"qbox.us/mgo2"
)

func ConfigDB(tea *teapot.Teapot) {
	log := tea.Logger()
	log.Debugf("db: %#v", global.Env.Mongo.Default)
	database, err := mgo2.NewDatabaseWithTimeoutNoFatal(global.Env.Mongo.Default, "strong", 2*time.Second)
	if err != nil {
		log.Errorf("mgo2.NewDatabaseWithTimeoutNoFatal(%s, strong) with error : %s", global.Env.Mongo.Default, err)
	}

	models.MongDbs[models.DefaultMgo] = func() *mgo2.Database {
		return database.Copy()
	}
}
