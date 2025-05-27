package middlewares

import (
	"log"
	"net/http"
//	"runtime/debug"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
//		log.Printf("%s", string(debug.Stack()))
		next.ServeHTTP(w, r)
	})
}
