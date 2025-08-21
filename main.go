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

	_ "github.com/golang-migrate/migrate/v4/database/postgres" // blank import
	_ "github.com/golang-migrate/migrate/v4/source/file"       // blank import

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
)

func runMigrations() {
	migrationsDir := os.Getenv("MIGRATIONS_DIR") // file://internal/db/migrations
	dbURL := os.Getenv("DATABASE_URL")           // postgres://...

	if migrationsDir == "" || dbURL == "" {
		log.Println("[migrate] MIGRATIONS_DIR veya DATABASE_URL boş, migration atlandı")
		return
	}

	m, err := migrate.New(migrationsDir, dbURL)
	if err != nil {
		log.Fatalf("[migrate] migrate.New hata: %v", err)
	}

	defer func() {
		src, db := m.Close()
		if src != nil {
			log.Printf("[migrate] source close: %v", src)
		}
		if db != nil {
			log.Printf("[migrate] database close: %v", db)
		}
	}()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("[migrate] up-to-date ✅ (uygulanacak değişiklik yok)")
		} else {
			log.Fatalf("[migrate] Up hata: %v", err)
		}
	} else {
		log.Println("[migrate] başarıyla uygulandı ✅")
	}
}

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

	// 1) Migration'ı önce çalıştır
	runMigrations()

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
	profilService := authsvc.NewServiceProfile(userRepo)
	profilHandler := auth.NewHandlerGetProfile(profilService)
	authHandler := auth.NewHandler(authService, jwtSvc)

	r := router.New(
		middleware.CORS,
		middleware.JWTAuth(jwtSvc), // protected rotalar için
		authHandler,
		profilHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
