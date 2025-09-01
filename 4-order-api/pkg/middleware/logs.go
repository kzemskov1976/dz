package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		log.SetFormatter(&log.JSONFormatter{})
		log.WithFields(log.Fields{
			"status_code": wrapper.StatusCode,
			"method":      r.Method,
			"path":        r.URL.Path,
		}).Info("")
	})
}
