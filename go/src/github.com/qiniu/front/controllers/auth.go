package controllers

import (
	"bytes"
	"encoding/json"

	"github.com/teapots/render"
	"hd.qiniu.com/models"
)

type Auth struct {
	Base
	Render render.Render `inject`
}

func (c *Auth) Get() {
	var (
		res = map[string]interface{}{
			"error":    nil,
			"logined":  false,
			"email":    "",
			"fullname": "",
		}
		msg   string
		state = c.Params.Query.Get("state")
		code  = c.Params.Query.Get("code")
	)

	defer func() {
		switch state {
		case "js":
			if msg != "" {
				res["error"] = msg
			}
			buf := bytes.NewBuffer([]byte{})

			c.Render.Header().Set("Content-Type", "text/javascript")
			buf.WriteString("auth_result_callback(\n")
			json.NewEncoder(buf).Encode(res)
			buf.WriteString(");")

			c.Render.Data(buf.Bytes())
		default:
			if msg == "" {
				msg = "登录成功"
			} else {
				res["error"] = msg
			}
			c.FlashMsg(msg)
			c.Render.HTML("auth", res)
		}
	}()

	token, err := c.ExchangeOAuthToken(code)
	if err != nil {
		c.Log.Errorf("ExchangeOAuthToken error: %s", err)
		msg = "获取授权出错"
		return
	}

	if token.Error != 0 {
		c.Log.Errorf("ExchangeOAuthToken return code: %d, error: %s", token.Error, token.ErrorDescription)
		msg = "获取授权出错"
		return
	}

	// 获取用户信息
	info, err := c.AccountService.GetAccountInfo(token.AccessToken)
	if err != nil {
		msg = "获取用户信息出错，请稍后再试"
		c.Log.Errorf("AccountService.GetAccountInfo(%s) error: %s", token.AccessToken, err)
		return
	}

	if info.Ret != 0 {
		c.Log.Errorf("AccountService.GetAccountInfo return %d, msg: %s", info.Ret, info.Msg)
		msg = "获取portal用户信息出错"
		return
	}

	c.Log.Debugf("info.Account is: %#v", info.Account)

	// 查找该用户注册信息
	user, err := models.User.FindByUid(info.Account.Uid)
	if err != nil {
		if err != models.ErrNotFound {
			c.Log.Errorf("models.User.FindByUid(%d) error: %s", info.Account.Uid, err)
			msg = "获取用户信息出错"
			return
		}

		c.Log.Debug("user is not found, create!")
		user = models.User.New(info.Account.Uid, info.Account.Email, info.Account.FullName, info.Account.Gender)
		err = models.User.Insert(user)
		if err != nil {
			c.Log.Errorf("models.User.Insert() error: %s", err)
			msg = "初始化用户出错"
			return
		}
	}

	c.Log.Debugf("user is: %#v", user)

	c.SetLoginUid(info.Account.Uid)

	// 持久化token
	err = c.SetOAuthToken(info.Account.Uid, token)
	if err != nil {
		c.Log.Error(err)
		msg = "保存登录状态出错"
		return
	}

	res = map[string]interface{}{
		"error":    nil,
		"logined":  true,
		"email":    user.Email,
		"fullname": user.FullName,
	}
	return
}
