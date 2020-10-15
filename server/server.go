package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Config struct {
	TempDir string `default:"/tmp"`
}

var conf Config

func NewRouter(config Config) *chi.Mux {

	conf = config

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
	r.Route("/report", func(r chi.Router) {
		r.Get("/account-transactions", reportAccountTransactionsHandler)
	})

	return r
}
