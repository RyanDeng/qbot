package mgo

import (
	"fmt"
	"labix.org/v2/mgo"
	"log"
	"strings"
	"time"
)

var (
	CopySessionMaxRetry = 5
)

// ------------------------------------------------------------------------

func Dail(host string, mode string, syncTimeoutInS int64) *mgo.Session {

	session, err := mgo.Dial(host)
	if err != nil {
		log.Fatal("Connect MongoDB failed:", err, host)
	}

	if mode != "" {
		SetMode(session, mode, true)
	}
	if syncTimeoutInS != 0 {
		session.SetSyncTimeout(time.Duration(int64(time.Second) * syncTimeoutInS))
	}

	return session
}

// ------------------------------------------------------------------------

type Config struct {
	Host           string `json:"host"`
	DB             string `json:"db"`
	Coll           string `json:"coll"`
	Mode           string `json:"mode"`
	SyncTimeoutInS int64  `json:"timeout"` // 以秒为单位
}

type Session struct {
	*mgo.Session
	DB   *mgo.Database
	Coll *mgo.Collection
}

func Open(cfg *Config) *Session {

	session := Dail(cfg.Host, cfg.Mode, cfg.SyncTimeoutInS)
	db := session.DB(cfg.DB)
	c := db.C(cfg.Coll)

	return &Session{session, db, c}
}

// test whether session closed
//
// PS: sometimes it's not corrected
func IsSessionClosed(s *mgo.Session) (res bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("[MGO2_IS_SESSION_CLOSED] check session closed panic:", err)
		}
	}()
	res = true
	return s.Ping() != nil
}

func checkSession(s *mgo.Session) (err error) {
	return s.Ping()
}

func isServersFailed(err error) bool {
	return strings.Contains(err.Error(), "no reachable servers")
}

func CopySession(s *mgo.Session) *mgo.Session {
	for i := 0; i < CopySessionMaxRetry; i++ {
		res := s.Copy()
		err := checkSession(res)
		if err == nil {
			return res
		}
		CloseSession(res)
		log.Print("[MGO2_COPY_SESSION] copy session and check failed:", err)
		if isServersFailed(err) {
			panic("[MGO2_COPY_SESSION_FAILED] servers failed")
		}
	}
	msg := fmt.Sprintf("[MGO2_COPY_SESSION_FAILED] failed after %d retries", CopySessionMaxRetry)
	log.Fatal(msg)
	panic(msg)
}

func FastCopySession(s *mgo.Session) *mgo.Session {
	return s.Copy()
}

func CloseSession(s *mgo.Session) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("[MGO2_CLOSE_SESSION_RECOVER] close session panic", err)
		}
	}()
	s.Close()
}

// copy database's session, and re-create the database.
//
// you need call `CloseDatabase` after use this
func CopyDatabase(db *mgo.Database) *mgo.Database {
	return CopySession(db.Session).DB(db.Name)
}

func FastCopyDatabase(db *mgo.Database) *mgo.Database {
	return FastCopySession(db.Session).DB(db.Name)
}

// close the session of the database
func CloseDatbase(db *mgo.Database) {
	CloseSession(db.Session)
}

// copy collection's session, and re-create the collection
//
// you need call `CloseColletion` after use this
func CopyCollection(c *mgo.Collection) *mgo.Collection {
	return CopyDatabase(c.Database).C(c.Name)
}

func FastCopyCollection(c *mgo.Collection) *mgo.Collection {
	return FastCopyDatabase(c.Database).C(c.Name)
}

// close the session of the collection
func CloseCollection(c *mgo.Collection) {
	CloseDatbase(c.Database)
}

func CheckIndex(c *mgo.Collection, key []string, unique bool) error {
	originIndexs, err := c.Indexes()
	if err != nil {
		return fmt.Errorf("<CheckIndex> get indexes: %v", err)
	}
	for _, index := range originIndexs {
		if checkIndexKey(index.Key, key) && unique == index.Unique {
			return nil
		}
	}
	return fmt.Errorf("<CheckIndex> not found index: %v unique: %v", key, unique)
}

func checkIndexKey(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, k := range a {
		if k != b[i] {
			return false
		}
	}
	return true
}

//------------------------------------
var g_modes = map[string]int{
	"eventual":  0,
	"monotonic": 1,
	"mono":      1,
	"strong":    2,
}

func SetMode(s *mgo.Session, modeFriendly string, refresh bool) {

	mode, ok := g_modes[strings.ToLower(modeFriendly)]
	if !ok {
		log.Fatal("invalid mgo mode")
	}
	switch mode {
	case 0:
		s.SetMode(mgo.Eventual, refresh)
	case 1:
		s.SetMode(mgo.Monotonic, refresh)
	case 2:
		s.SetMode(mgo.Strong, refresh)
	}
}
