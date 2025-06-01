package middleware

import (
	"context"
	"javaneseivankov/url-shortener/internal/contextkeys"
	"javaneseivankov/url-shortener/internal/errx"
	"javaneseivankov/url-shortener/pkg"
	"javaneseivankov/url-shortener/pkg/jwt"
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
				pkg.SendError(w, errx.ErrBearerTokenInvalidFormat)
				return
			}

			tokenString = parts[1]
			claims, err := jwt.VerifyToken(tokenString)
			if err != nil {
				pkg.SendError(w, errx.ErrInvalidBearerToken)
				return
			}

			// r.Header.Set("user_id", claims.UserID.String())
			// r.Header.Set("username", claims.Username)
			// r.Header.Set("role", claims.Role)

			ctx := context.WithValue(r.Context(), contextkeys.ClaimCtxKey, claims)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			http.Error(w, "Unathorized: No token provided", http.StatusUnauthorized)
// 			return
// 		}

// 		parts := strings.Split(tokenString, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString = parts[1]
// 		var claims pkg.Claims
// 		err := (tokenString, &claims)
// 		if err != nil {
// 			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		r.Header.Set("user_id", claims.UserID)
// 		r.Header.Set("username", claims.Username)
// 		r.Header.Set("role", claims.Role)

// 		next.ServeHTTP(w, r)
// 	})
// }
