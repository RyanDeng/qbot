package session

import (
	"net/http"

	"github.com/teapots/teapot"
	"hd.qiniu.com/utils/sessions"
)

func SessionStore() interface{} {
	return func(log teapot.Logger, req *http.Request, rw http.ResponseWriter, manager *sessions.SessionManager) sessions.SessionStore {
		sess, _, err := manager.Start(rw, req)
		if err != nil {
			log.Error(err)
			return nil
		}

		if trw, ok := rw.(teapot.ResponseWriter); ok {
			trw.Before(func(rw teapot.ResponseWriter) {
				err = sess.Flush()
				if err != nil {
					log.Warn("sess flush:", err)
				}
			})
		}

		return sess
	}
}
