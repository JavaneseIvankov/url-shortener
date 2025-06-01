package middleware

import "net/http"

func ApplyMiddleware(f http.HandlerFunc, mws ...func(http.Handler) http.Handler) http.HandlerFunc {
    h := http.Handler(f)
    for _, mw := range mws {
        h = mw(h)
    }
    return handlerToFunc(h)
}

func handlerToFunc(h http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
