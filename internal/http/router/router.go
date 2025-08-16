package router

import (
	"net/http"

	"myfirstbackend/internal/user"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// New uygulamanın tüm HTTP rotalarını ve middleware'lerini kurar.
// CORS gibi global middleware'leri dışarıdan alır (test edilebilirlik için).
func New(cors func(http.Handler) http.Handler, uh *user.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	if cors != nil {
		r.Use(cors)
	}

	// Sağlık kontrolü
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// Swagger UI (/swagger/index.html)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// /api altında çalışan alt router
	api := chi.NewRouter()

	// ------- AUTH -------
	api.Route("/auth", func(r chi.Router) {
		r.Post("/register", uh.Register) // POST /api/auth/register
		// ileride: r.Post("/login", uh.Login) vb.
	})

	// Gerekirse burada başka alanları mount edebilirsin
	// api.Mount("/users", usersRouter(uh)) vb.

	// /api'ye mount et
	r.Mount("/api", api)

	return r
}
