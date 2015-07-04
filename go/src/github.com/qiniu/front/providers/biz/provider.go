package biz

import (
	"net/http"

	"hd.qiniu.com/services/biz"
)

func BizService(bizHost string, tr http.RoundTripper) biz.BizService {
	return biz.NewBizService(bizHost, tr)
}
