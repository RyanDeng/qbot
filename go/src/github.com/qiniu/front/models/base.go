package models

import "qbox.us/mgo2"

const (
	DefaultMgo = "default"
)

var (
	MongDbs = make(map[string]func() *mgo2.Database)
)
