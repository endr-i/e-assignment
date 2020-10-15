package repo

import (
	"assignment/repo/account"
	"assignment/repo/operation"
	rateRepository "assignment/repo/rate"
	registerRepository "assignment/repo/register"
	reportRepository "assignment/repo/report"
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once
)

func Init(db *gorm.DB) {
	// could create repos for different models if multi database
	once.Do(func() {
		registerRepository.Init(db)
		accountRepo := accountRepository.InitRepo(db)
		rateRepo := rateRepository.InitRepo(db)
		operationRepo := operationRepository.InitRepo(db, accountRepo, rateRepo)
		reportRepository.InitRepo(db, accountRepo, operationRepo)
	})
}
