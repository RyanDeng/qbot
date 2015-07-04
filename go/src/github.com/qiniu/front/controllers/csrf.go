package controllers

import (
	"net/http"

	"github.com/teapots/teapot"
	"hd.qiniu.com/utils/api"
)

func CSRFHandler(ctx teapot.Context, rw http.ResponseWriter, req *http.Request, log teapot.Logger) {
	log.Notice(api.CSRFDetected.Humanize())

	(&api.JsonResult{
		Code: api.CSRFDetected,
	}).Write(ctx, rw, req)
}
