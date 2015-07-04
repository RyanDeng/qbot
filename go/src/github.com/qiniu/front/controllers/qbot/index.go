package qbot

import (
	"github.com/teapots/render"
	"hd.qiniu.com/controllers"
)

type Index struct {
	controllers.Base
	render.Render `inject`
}

func (c *Index) Get() {
	c.HTML("qbot/index", nil)
}
