package api

import (
	"encoding/json"
	"net/http"

	"github.com/teapots/teapot"
)

var JsonContentType = "application/json; charset=UTF-8"

type JsonResult struct {
	Code    Code        `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func (r *JsonResult) Write(ctx teapot.Context, rw http.ResponseWriter, req *http.Request) {
	if r.Code == 0 {
		r.Code = OK
	}

	if r.Message == "" && r.Code != OK {
		r.Message = r.Code.Humanize()
	}

	var config *teapot.Config
	ctx.Find(&config, "")

	var body []byte
	var err error
	if config.RunMode.IsDev() {
		body, err = json.MarshalIndent(r, "", "  ")
	} else {
		body, err = json.Marshal(r)
	}

	rw.Header().Set("Content-Type", JsonContentType)

	if err == nil {
		_, err = rw.Write(body)
	} else {
		v := &JsonResult{
			Code:    ResultError,
			Message: ResultError.Humanize(),
		}
		if config.RunMode.IsDev() {
			body, _ = json.MarshalIndent(v, "", "  ")
		} else {
			body, _ = json.Marshal(v)
		}
		rw.Write(body)
	}

	if err != nil {
		var logger teapot.Logger
		ctx.Find(&logger, "")
		logger.Error("JsonResult.Write", err)
	}
}
