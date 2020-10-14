package repo

import (
	"assignment/repo/account"
	"assignment/repo/operation"
	rateRepository "assignment/repo/rate"
	registerRepository "assignment/repo/register"
	"errors"
	"gorm.io/gorm"
	"sync"
)

var (
	once             sync.Once
	InvalidUuidError = errors.New("invalid uuid")
	NoCurrencyError  = errors.New("no currency")
	NoUserError      = errors.New("no user")
	NoAccountError   = errors.New("no account")
	NoRateError      = errors.New("no rate")
	LowBalance       = errors.New("low balance")
	NoRatesToUpload  = errors.New("no rates to upload")
)

func Init(db *gorm.DB) {
	// could create repos for different models if multi database
	once.Do(func() {
		registerRepository.Init(db)
		accountRepo := accountRepository.InitRepo(db)
		rateRepo := rateRepository.InitRepo(db)
		operationRepository.InitRepo(db, accountRepo, rateRepo)
	})
}
