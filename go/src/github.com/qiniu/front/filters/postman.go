package filters

import (
	"net/http"
	"strings"
)

func HeaderRemovePrefixFilter(prefix string) interface{} {
	return func(req *http.Request) {
		for key, value := range req.Header {
			if !strings.HasPrefix(key, prefix) {
				continue
			}
			name := strings.TrimPrefix(key, prefix)
			if _, ok := req.Header[name]; ok {
				continue
			}
			req.Header[name] = value
			delete(req.Header, key)
		}
	}
}
