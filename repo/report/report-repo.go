package reportRepository

import (
	accountRepository "assignment/repo/account"
	operationRepository "assignment/repo/operation"
	"assignment/utils"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	repo *repository
	once sync.Once
)

type IRepository interface {
	AccountTransactionsReport(AccountTransactionsReportForm) (*AccountTransactionsReportData, error)
}

type repository struct {
	db            *gorm.DB
	accountRepo   accountRepository.IRepository
	operationRepo operationRepository.IRepository
}

func InitRepo(db *gorm.DB, accountRepo accountRepository.IRepository, operationRepo operationRepository.IRepository) IRepository {
	once.Do(func() {
		repo = &repository{
			db:            db,
			accountRepo:   accountRepo,
			operationRepo: operationRepo,
		}
	})
	return repo
}

func GetRepo() IRepository {
	return repo
}

func (r *repository) AccountTransactionsReport(form AccountTransactionsReportForm) (*AccountTransactionsReportData, error) {
	now := time.Now()

	since, err := time.ParseInLocation("2006-01-02", form.Date, time.Local)

	if err != nil {
		return nil, utils.InvalidDateFormatError
	}
	till := since.Add(24*time.Hour - time.Second) // till next day

	account, err := r.accountRepo.GetById(form.AccountId)
	if err != nil {
		return nil, err
	}

	transactions, err := r.operationRepo.AccountTransactions(account.ID, since, till)
	if err != nil {
		return nil, err
	}

	reportTransactions := make([]AccountTransactionsReportTransaction, 0, len(transactions))
	for _, transaction := range transactions {
		reportTransactions = append(reportTransactions, GetAccountTransactionsReportTransaction(&transaction))
	}

	return &AccountTransactionsReportData{
		Account:      account,
		Transactions: reportTransactions,
		Since:        since,
		Till:         till,
		DateTime:     now,
	}, nil
}
