package main

import (
	"log"
	"net/http"
	"os"

	"myfirstbackend/internal/db"
	mw "myfirstbackend/internal/http/middleware"
	"myfirstbackend/internal/user"

	_ "myfirstbackend/docs" // Import the generated docs package

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
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

	repo := user.NewRepository(pool)
	uh := user.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(mw.CORS)

	// Swagger UI ( /swagger/index.html )
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Auth routes
	r.Post("/api/auth/register", uh.Register)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
