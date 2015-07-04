package qbot

import (
	"github.com/teapots/render"
	"hd.qiniu.com/controllers"
	"strconv"
	"strings"

	"hd.qiniu.com/models/qbot"
	"net/http"
)

type Admin struct {
	controllers.Base
	render.Render `inject`
}

func (c *Admin) Post() {
	var (
		res = map[string]interface{}{
			"message": "添加成功",
		}
		statusCode = http.StatusOK
		Name       = strings.TrimSpace(c.Params.Get("name"))
		Phone      = strings.TrimSpace(c.Params.Get("phone"))
		Department = strings.TrimSpace(c.Params.Get("department"))
		Email      = strings.TrimSpace(c.Params.Get("email"))
		NickName   = strings.TrimSpace(c.Params.Get("nickname"))
		QQ         = strings.TrimSpace(c.Params.Get("qq"))
	)
	phoneNum, _ := strconv.ParseInt(Phone, 10, 64)
	defer func() {
		res["code"] = statusCode
		c.Render.JSON(res)
	}()
	infor := qbot.Qbot.New(Name, NickName, Email, QQ, Department, phoneNum)

	err := qbot.Qbot.Insert(infor)
	if err != nil {
		c.Log.Errorf("create infor error: %s", err)
		statusCode = http.StatusServiceUnavailable
		res["message"] = "添加出错，请稍后再试"
	}
}
