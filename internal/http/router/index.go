package router

import (
	"net/http"

	"myfirstbackend/internal/handler/auth"
	//"myfirstbackend/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

func New(cors func(http.Handler) http.Handler, jwtMw func(http.Handler) http.Handler, ah *auth.Handler) *chi.Mux {
	r := chi.NewRouter()
	if cors != nil {
		r.Use(cors)
	}

	// /health
	MountHealth(r)

	// /swagger/*
	MountSwagger(r)

	// /api/*
	api := chi.NewRouter()
	MountAuth(api, ah) // /api/auth/*
	// MountUser(api, uh)        // ileride: /api/users/*
	// MountWhatever(api, ...)

	// Protected örneği (ileride):
	// api.With(jwtMw).Route("/profile", func(pr chi.Router) {
	//   pr.Get("/", profileHandler)
	// })

	r.Mount("/api", api)
	return r
}
