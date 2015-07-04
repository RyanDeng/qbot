package filters

import (
	"net/http"
	"strings"
	"time"
)

func DurationFilter(startTime, endTime string) interface{} {
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		start = time.Now()
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		end = time.Now()
	}
	return func(req *http.Request, rw http.ResponseWriter) {
		var (
			now      = time.Now()
			pathInfo = strings.Split(req.URL.Path, "/")
			activity = ""
			action   = ""
		)

		if len(pathInfo) > 1 {
			activity = pathInfo[1]
		}

		if len(pathInfo) > 2 {
			action = pathInfo[2]
		}

		if action == "prelude" || action == "ending" {
			if now.Before(end) && now.After(start) {
				rw.WriteHeader(302)
				rw.Header().Set("Location", "/"+activity+"/")
				return
			}
		} else {
			if now.Before(start) {
				rw.WriteHeader(302)
				rw.Header().Set("Location", "/"+activity+"/prelude")
				return
			}
			if now.After(end) {
				rw.WriteHeader(302)
				rw.Header().Set("Location", "/"+activity+"/ending")
			}
		}
	}
}
