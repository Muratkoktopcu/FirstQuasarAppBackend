package router

import (
	"net/http"

	"myfirstbackend/internal/handler/auth"

	"github.com/go-chi/chi/v5"
)

func MountProfile(r chi.Router, jwtMw func(http.Handler) http.Handler, h *auth.HandlerGetProfile) {
	r.With(jwtMw).Route("/profile", func(pr chi.Router) {
		pr.Get("/", h.Me)
	})
}
