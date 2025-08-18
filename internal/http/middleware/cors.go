package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

func CORS(next http.Handler) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	allowedMethods := strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ",")
	allowedHeaders := strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ",")
	allowCredentials := os.Getenv("CORS_ALLOW_CREDENTIALS") == "true"

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: allowCredentials,
	})
	return c.Handler(next)
}
