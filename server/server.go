package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Config struct {
	Port string `default:"3000"`
}

func NewRouter(config Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/reg", registerHandler)
	r.Route("/operation", func(r chi.Router) {
		r.Post("/refill", operationRefillHandler)
		r.Post("/transfer", operationTransferHandler)
	})
	r.Route("/rate", func(r chi.Router) {
		r.Post("/upload", rateUploadHandler)
	})

	return r
}
