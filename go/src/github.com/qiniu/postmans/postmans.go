package postmans

import (
	. "github.com/qiniu/postmans/interfaces"
	_ "github.com/qiniu/postmans/qq"
)

func Get(name string) (func(string) (Postman, error), bool) {

	return GetPostG(name)
}
