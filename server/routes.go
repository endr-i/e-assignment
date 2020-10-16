package server

import (
	"assignment/repo/operation"
	rateRepository "assignment/repo/rate"
	"assignment/repo/register"
	reportRepository "assignment/repo/report"
	"assignment/utils"
	"fmt"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

// @Summary Creates new user + account
// @ID register
// @Accept  json
// @Produce  json
// @Param website body registerRepository.Form true "params"
// @Success 200 {object} entity.Account
// @Router /reg [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithTime(time.Now())
	var form registerRepository.Form
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		renderError(w, r, "cannot parse JSON", http.StatusBadRequest)
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
	logger.WithField("request", form)
	registerRepo := registerRepository.GetRepo()
	account, err := registerRepo.Register(form)
	if err != nil {
		logger.Error(err.Error())
		renderError(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info(account)
	render.JSON(w, r, account)
}

// @Summary Refill operation
// @ID operationRefill
// @Accept  json
// @Produce  json
// @Param website body operationRepository.RefillForm true "params"
// @Success 200 {object} entity.Operation
// @Router /operation/refill [post]
func operationRefillHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithTime(time.Now())
	var form operationRepository.RefillForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		renderError(w, r, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	logger.WithField("request", form)
	operationRepo := operationRepository.GetRepo()
	operation, err := operationRepo.Refill(form)
	if err != nil {
		logger.Error(err.Error())
		renderError(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info(operation)
	render.JSON(w, r, operation)
}

// @Summary Transfer operation
// @ID operationTransfer
// @Accept  json
// @Produce  json
// @Param website body operationRepository.TransferForm true "params"
// @Success 200 {object} entity.Operation
// @Router /operation/transfer [post]
func operationTransferHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithTime(time.Now())
	var form operationRepository.TransferForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		renderError(w, r, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	logger.WithField("request", form)
	operationRepo := operationRepository.GetRepo()
	operation, err := operationRepo.Transfer(form)
	if err != nil {
		logger.Error(err.Error())
		renderError(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info(operation)
	render.JSON(w, r, operation)
}

// @Summary Upload new rates
// @ID rateUpload
// @Accept  json
// @Produce  json
// @Param website body rateRepository.UploadRatesForm true "params"
// @Success 200 {array} entity.Rate
// @Router /rate/upload [post]
func rateUploadHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithTime(time.Now())
	var form rateRepository.UploadRatesForm
	err := render.DecodeJSON(r.Body, &form)
	if err != nil {
		renderError(w, r, "cannot parse JSON", http.StatusBadRequest)
		return
	}
	logger.WithField("request", form)
	rateRepo := rateRepository.GetRepo()
	rates, err := rateRepo.UploadRates(form)

	if err != nil {
		logger.Error(err.Error())
		renderError(w, r, err.Error(), http.StatusBadRequest)
	}
	logger.Info(rates)
	render.JSON(w, r, rates)
}

// @Summary Generates report
// @ID reportAccountTransactions
// @Produce  json
// @Param date query string true "Date of report (eg 2020-10-16)"
// @Param accountId query string true "Account id"
// @Param file query string false "file output"
// @Success 200 {object} reportRepository.AccountTransactionsReportData
// @Router /report/account-transactions [get]
func reportAccountTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.WithTime(time.Now())
	form := reportRepository.AccountTransactionsReportForm{
		Date:      r.URL.Query().Get("date"),
		AccountId: r.URL.Query().Get("accountId"),
	}

	fileResponse := len(r.URL.Query().Get("file")) > 0

	logger.WithField("request", form)

	reportRepo := reportRepository.GetRepo()
	reportData, err := reportRepo.AccountTransactionsReport(form)
	if err != nil {
		logger.Error(err.Error())
		renderError(w, r, err.Error(), http.StatusBadRequest)
	}
	logger.Info(reportData)
	if fileResponse {
		fileName := fmt.Sprintf("%s%s.csv", "AccountTransactionsReport", time.Now().String())
		filePath := fmt.Sprintf("%s/%s", os.TempDir(), fileName)
		err := utils.CreateCSV(reportData.GetCSVData(), filePath)
		if err != nil {
			logger.Error(err.Error())
			renderError(w, r, err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
			http.ServeFile(w, r, filePath)
		}
	} else {
		render.JSON(w, r, reportData)
	}
}
