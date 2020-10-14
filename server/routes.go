package server

import (
	"assignment/repo"
	"assignment/repo/operation"
	rateRepository "assignment/repo/rate"
	"assignment/repo/register"
	"errors"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var form registerRepository.Form
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
	registerRepo := registerRepository.GetRepo()
	account, err := registerRepo.Register(form)
	if err != nil {
		if errors.Is(err, repo.NoCurrencyError) {
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

func operationRefillHandler(w http.ResponseWriter, r *http.Request) {
	var form operationRepository.RefillForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		http.Error(w, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	operationRepo := operationRepository.GetRepo()
	operation, err := operationRepo.Refill(form)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	render.JSON(w, r, operation)
}

func operationTransferHandler(w http.ResponseWriter, r *http.Request) {
	var form operationRepository.TransferForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		http.Error(w, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	operationRepo := operationRepository.GetRepo()
	operation, err := operationRepo.Transfer(form)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	render.JSON(w, r, operation)
}

func rateUploadHandler(w http.ResponseWriter, r *http.Request) {
	var form rateRepository.UploadRatesForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		http.Error(w, "cannot parse JSON", http.StatusBadRequest)
		return
	}

	rateRepo := rateRepository.GetRepo()
	rates, err := rateRepo.UploadRates(form)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	render.JSON(w, r, rates)
}
