package router

import (
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

func MountSwagger(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
