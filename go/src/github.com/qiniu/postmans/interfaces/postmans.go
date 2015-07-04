package interfaces

import (
//	. "github.com/qiniu/postmans/interfaces"
)

var (
	postmans = make(map[string]func(string) (Postman, error))
)

func Register(name string, manGen func(string) (Postman, error)) {
	postmans[name] = manGen
}

func GetPostG(name string) (man func(string) (Postman, error), ok bool) {
	man, ok = postmans[name]
	return
}
