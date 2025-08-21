package router

import (
	"net/http"

	"myfirstbackend/internal/handler/auth"
	//"myfirstbackend/internal/http/middleware"

	"github.com/go-chi/chi/v5"
)

// yeni bir router oluşturur(chi.Mux) ve CORS middleware'ini uygular.
func New(cors func(http.Handler) http.Handler, jwtMw func(http.Handler) http.Handler, ah *auth.Handler, ph *auth.HandlerGetProfile) *chi.Mux {
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
	MountProfile(api, jwtMw, ph)

	r.Mount("/api", api) // /api/* rotalarını api router'ına bağlar
	/*Mount → Başka bir router veya handler’ı, belli bir path’in altına takmaya (mount etmeye) yarar.
	  Daha büyük projelerde endpointleri ayırmak için çok kullanılır.*/
	return r
}
