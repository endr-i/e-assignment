package operationRepository

import (
	"assignment/entity"
	accountRepository "assignment/repo/account"
	rateRepository "assignment/repo/rate"
	"assignment/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	repo *repository
	once sync.Once
)

type IRepository interface {
	Refill(RefillForm) (*entity.Operation, error)
	Transfer(TransferForm) (*entity.Operation, error)
	AccountTransactions(accountId uuid.UUID, since time.Time, till time.Time) ([]entity.Transaction, error)
}

type repository struct {
	db          *gorm.DB
	accountRepo accountRepository.IRepository
	rateRepo    rateRepository.IRepository
}

func InitRepo(db *gorm.DB, accountRepo accountRepository.IRepository, rateRepo rateRepository.IRepository) IRepository {
	once.Do(func() {
		repo = &repository{
			db:          db,
			accountRepo: accountRepo,
			rateRepo:    rateRepo,
		}
	})
	return repo
}

func GetRepo() IRepository {
	return repo
}

func (r *repository) Refill(form RefillForm) (*entity.Operation, error) {
	account, err := r.accountRepo.GetById(form.AccountId)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	var currency entity.Currency
	if err := r.db.First(&currency, "symbol=?", form.Currency).Error; err != nil {
		return nil, utils.NoCurrencyError
	}

	currencyRate, err := r.rateRepo.GetCurrencyRate(currency.ID, now)
	if err != nil {
		return nil, err
	}
	accountRate, err := r.rateRepo.GetCurrencyRate(account.CurrencyId, now)
	if err != nil {
		return nil, err
	}

	value := decimal.NewFromFloat(form.Sum)

	operationDetails := entity.OperationRefillDetails{Source: form.Details.Source}
	operation := entity.Operation{
		ID:       uuid.New(),
		Type:     entity.OperationTypeRefill,
		Details:  operationDetails.JSON(),
		DateTime: now,
	}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&operation).Error; err != nil {
			return err
		}
		fillTransaction := entity.Transaction{
			ID:                uuid.New(),
			OperationId:       operation.ID,
			Value:             value,
			AccountId:         account.ID,
			Account:           account,
			DateTime:          now,
			CurrencyId:        currency.ID,
			Currency:          &currency,
			CurrencyRateValue: currencyRate.Value,
			AccountRateValue:  accountRate.Value,
		}
		if err = tx.Create(&fillTransaction).Error; err != nil {
			return err
		}
		operation.Transactions = []entity.Transaction{fillTransaction}

		diff := fillTransaction.RatedValue()
		err = tx.Model(account).Where("id=?", account.ID).
			Update("balance", account.Balance.Add(diff)).
			Scan(account).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operation, err
}

func (r *repository) Transfer(form TransferForm) (*entity.Operation, error) {
	accountFrom, err := r.accountRepo.GetById(form.From)
	if err != nil {
		return nil, err
	}

	accountTo, err := r.accountRepo.GetById(form.To)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	rateFrom, err := r.rateRepo.GetCurrencyRate(accountFrom.CurrencyId, now)
	if err != nil {
		return nil, err
	}
	rateTo, err := r.rateRepo.GetCurrencyRate(accountTo.CurrencyId, now)
	if err != nil {
		return nil, err
	}

	var (
		currency     *entity.Currency
		currencyRate *entity.Rate
	)

	if form.Currency == accountTo.Currency.Symbol {
		currency = accountTo.Currency
		currencyRate = rateTo
	} else if form.Currency == accountFrom.Currency.Symbol {
		currency = accountFrom.Currency
		currencyRate = rateFrom
	} else {
		return nil, utils.ForbiddenCurrencyError
	}

	value := decimal.NewFromFloat(form.Sum)

	operationDetails := entity.OperationTransferDetails{Comment: form.Details.Comment}
	operation := entity.Operation{
		ID:       uuid.New(),
		Type:     entity.OperationTypeTransfer,
		Details:  operationDetails.JSON(),
		DateTime: now,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&operation).Error; err != nil {
			return err
		}
		debitTransaction := entity.Transaction{
			ID:                uuid.New(),
			OperationId:       operation.ID,
			Value:             value.Neg(),
			AccountId:         accountFrom.ID,
			Account:           accountFrom,
			DateTime:          now,
			CurrencyId:        currency.ID,
			Currency:          currency,
			CurrencyRateValue: currencyRate.Value,
			AccountRateValue:  rateFrom.Value,
		}

		if accountFrom.Balance.Add(debitTransaction.RatedValue()).IsNegative() {
			return utils.LowBalanceError
		}

		if err := tx.Create(&debitTransaction).Error; err != nil {
			return err
		}

		refillTransaction := entity.Transaction{
			ID:                uuid.New(),
			OperationId:       operation.ID,
			Value:             value,
			AccountId:         accountTo.ID,
			Account:           accountTo,
			DateTime:          now,
			CurrencyId:        currency.ID,
			Currency:          currency,
			CurrencyRateValue: currencyRate.Value,
			AccountRateValue:  rateTo.Value,
		}
		if err := tx.Create(&refillTransaction).Error; err != nil {
			return err
		}

		operation.Transactions = []entity.Transaction{debitTransaction, refillTransaction}

		err := tx.Model(accountFrom).Where("id=?", accountFrom.ID).
			Update("balance", accountFrom.Balance.Add(debitTransaction.RatedValue())).
			Scan(accountFrom).Error
		if err != nil {
			return err
		}

		err = tx.Model(accountTo).Where("id=?", accountTo.ID).
			Update("balance", accountTo.Balance.Add(refillTransaction.RatedValue())).
			Scan(accountTo).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &operation, nil
}

func (r *repository) AccountTransactions(accountId uuid.UUID, since time.Time, till time.Time) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Joins("Operation").Joins("Currency").Where("account_id=? and transactions.date_time between ? and ?", accountId, since, till).Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
