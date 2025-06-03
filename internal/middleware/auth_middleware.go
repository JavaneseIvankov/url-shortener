package middleware

import (
	"context"
	"javaneseivankov/url-shortener/internal/contextkeys"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/jwt"
	"javaneseivankov/url-shortener/pkg/logger"
	"net/http"
	"strings"
)

func AuthMiddleware(jwt jwt.JWT) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				pkg.SendError(w, errx.ErrNoBearerToken)
			}
			

			parts := strings.Split(tokenString, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Debug("Bearer token splitted", "auth_header_value", tokenString,"partsLength", len(parts))
				pkg.SendError(w, errx.ErrBearerTokenInvalidFormat)
				return
			}

			tokenString = parts[1]
			claims, err := jwt.VerifyToken(tokenString)
			if err != nil {
				pkg.SendError(w, errx.ErrInvalidBearerToken)
				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.ClaimCtxKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
