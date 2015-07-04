package env

import (
	. "github.com/teapots/teapot"
	"hd.qiniu.com/controllers"
	"hd.qiniu.com/controllers/qbot"

	"hd.qiniu.com/filters"
)

func ConfigRoutes(tea *Teapot) {
	tea.Routers(

		Router("/qbot",
			Router("/post", Post(&qbot.Admin{})),
			Get(&qbot.Index{}),
		),
	)
}
