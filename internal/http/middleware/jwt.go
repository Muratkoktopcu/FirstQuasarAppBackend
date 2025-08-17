package middleware

import (
	"context"
	"net/http"
	"strings"

	"myfirstbackend/internal/security/jwt"
)

type ctxKey string

const ContextUserID ctxKey = "userID"

func JWTAuth(jwtSvc jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			_, claims, err := jwtSvc.Parse(tokenStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			sub, ok := claims["sub"].(float64)
			if !ok {
				http.Error(w, "invalid subject", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserID, int64(sub))
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
