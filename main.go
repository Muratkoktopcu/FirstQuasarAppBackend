package main

import (
	"log"
	"net/http"
	"os"

	"myfirstbackend/internal/db"
	"myfirstbackend/internal/handler/auth"
	"myfirstbackend/internal/http/middleware"
	"myfirstbackend/internal/http/router"
	"myfirstbackend/internal/repository/user"
	"myfirstbackend/internal/security/jwt"
	authsvc "myfirstbackend/internal/service/auth"

	_ "myfirstbackend/docs"

	"github.com/joho/godotenv"
)

// @title           My API
// @version         1.0
// @description     Quasar frontend'in kullandığı Go backend API dokümantasyonu.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load(".env", "../.env", "../../.env")

	pool := db.Connect()
	defer pool.Close()

	userRepo := user.NewPgRepository(pool)

	jwtCfg := jwt.Config{
		Secret:          os.Getenv("JWT_SECRET"),
		AccessTokenTTL:  os.Getenv("JWT_ACCESS_TTL"),  // örn: "15m"
		RefreshTokenTTL: os.Getenv("JWT_REFRESH_TTL"), // örn: "720h"
		Issuer:          "myfirstbackend",
	}
	jwtSvc := jwt.New(jwtCfg)

	authService := authsvc.New(userRepo, jwtSvc)
	authHandler := auth.NewHandler(authService)

	r := router.New(
		middleware.CORS,
		middleware.JWTAuth(jwtSvc), // protected rotalar için
		authHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
