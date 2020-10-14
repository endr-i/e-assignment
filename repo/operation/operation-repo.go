package operationRepository

import (
	"assignment/entity"
	repo2 "assignment/repo"
	accountRepository "assignment/repo/account"
	rateRepository "assignment/repo/rate"
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

type RefillForm struct {
	AccountId string
	Sum       float64
	Currency  string
	Details   struct {
		Source string
	} // Details of operation
}

type TransferForm struct {
	From     string
	To       string
	Sum      float64
	Currency string
	Details  struct {
		Comment string
	} // Details of operation
}

type IRepository interface {
	Refill(RefillForm) (*entity.Operation, error)
	Transfer(TransferForm) (*entity.Operation, error)
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
	accountUuid, err := uuid.Parse(form.AccountId)
	if err != nil {
		return nil, repo2.InvalidUuidError
	}
	account, err := r.accountRepo.GetById(accountUuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	currencyId, err := uuid.Parse(form.Currency)
	if err != nil {
		return nil, repo2.InvalidUuidError
	}

	convertRate, err := r.rateRepo.GetConvertRate(currencyId, account.CurrencyId, now)

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
			ID:          uuid.New(),
			OperationId: operation.ID,
			Value:       value,
			AccountId:   account.ID,
			Account:     account,
			DateTime:    now,
			CurrencyId:  currencyId,
			RateValue:   convertRate.RateValue,
		}
		if err = tx.Create(&fillTransaction).Error; err != nil {
			return err
		}
		operation.Transactions = []entity.Transaction{fillTransaction}

		diff := fillTransaction.DiffValue()
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
	accountFromUuid, err := uuid.Parse(form.From)
	if err != nil {
		return nil, repo2.InvalidUuidError
	}
	accountFrom, err := r.accountRepo.GetById(accountFromUuid)
	if err != nil {
		return nil, err
	}

	accountToUuid, err := uuid.Parse(form.From)
	if err != nil {
		return nil, repo2.InvalidUuidError
	}
	accountTo, err := r.accountRepo.GetById(accountToUuid)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	convertRate, err := r.rateRepo.GetConvertRate(accountFrom.CurrencyId, accountTo.CurrencyId, now)
	if err != nil {
		return nil, err
	}

	value := decimal.NewFromFloat(form.Sum)
	diff := value.Mul(convertRate.RateValue)

	balanceFrom := accountFrom.Balance.Sub(diff)
	if balanceFrom.IsNegative() {
		return nil, repo2.LowBalance
	}

	operationDetails := entity.OperationTransferDetails{Comment: form.Details.Comment}
	operation := entity.Operation{
		ID:       uuid.New(),
		Type:     entity.OperationTypeRefill,
		Details:  operationDetails.JSON(),
		DateTime: now,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&operation).Error; err != nil {
			return err
		}
		debitTransaction := entity.Transaction{
			ID:          uuid.New(),
			OperationId: operation.ID,
			Value:       value.Neg(),
			AccountId:   accountFrom.ID,
			Account:     accountFrom,
			DateTime:    now,
			CurrencyId:  accountFrom.CurrencyId,
			RateValue:   convertRate.RateValue,
		}
		if err := tx.Create(&debitTransaction).Error; err != nil {
			return err
		}
		refillTransaction := entity.Transaction{
			ID:          uuid.New(),
			OperationId: operation.ID,
			Value:       value,
			AccountId:   accountTo.ID,
			Account:     accountTo,
			DateTime:    now,
			CurrencyId:  accountTo.CurrencyId,
			RateValue:   convertRate.RateValue,
		}
		if err := tx.Create(&refillTransaction).Error; err != nil {
			return err
		}

		operation.Transactions = []entity.Transaction{debitTransaction, refillTransaction}

		err := tx.Model(accountFrom).Where("id=?", accountFrom.ID).
			Update("balance", accountFrom.Balance.Add(debitTransaction.DiffValue())).
			Scan(accountFrom).Error
		if err != nil {
			return err
		}

		err = tx.Model(accountFrom).Where("id=?", accountTo.ID).
			Update("balance", accountTo.Balance.Add(refillTransaction.DiffValue())).
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
