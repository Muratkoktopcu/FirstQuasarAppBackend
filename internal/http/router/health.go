package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Health godoc
// @Summary      Health check
// @Tags         Health
// @Success      200 {string} string "ok"
// @Router       /health [get]
func MountHealth(r *chi.Mux) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok")) // http responseun body'sine cevap yazma fonksiyonu
	})
}
