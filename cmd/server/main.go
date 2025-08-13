package main

import (
	"log"
	"net/http"
	"os"

	"myfirstbackend/internal/db"
	mw "myfirstbackend/internal/http/middleware"
	"myfirstbackend/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	pool := db.Connect()
	defer pool.Close()

	repo := user.NewRepository(pool)
	uh := user.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(mw.CORS)

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
