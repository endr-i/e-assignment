package server

import (
	"assignment/repo/register"
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var form register.Form
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		http.Error(w, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	//err := decodeJSONBody(w, r, &form)
	//if err != nil {
	//	var mr *malformedRequest
	//	if errors.As(err, &mr) {
	//		http.Error(w, mr.msg, mr.status)
	//	} else {
	//		log.Println(err.Error())
	//		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	}
	//	return
	//}
	registerRepo := register.GetRepo()
	account, err := registerRepo.Register(form)
	if err != nil {
		if errors.Is(err, register.NoCurrencyError) {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	render.JSON(w, r, account)
}
