package server

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Message string
}

func renderError(w http.ResponseWriter, r *http.Request, message string, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, ErrorResponse{Message: message})
}
