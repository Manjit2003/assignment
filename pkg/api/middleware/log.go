package middleware

import (
	"net/http"

	"github.com/Manjit2003/samespace/pkg/utils"
)

var log = utils.GetChildLogger("http")

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("request", "url", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
