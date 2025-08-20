package router

import (
	"myfirstbackend/internal/handler/auth"

	"github.com/go-chi/chi/v5"
)

func MountAuth(r chi.Router, h *auth.Handler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register) // POST /api/auth/register
		r.Post("/login", h.Login)       // POST /api/auth/login
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
	})
}
