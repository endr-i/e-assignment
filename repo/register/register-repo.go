package registerRepository

import (
	"assignment/entity"
	"assignment/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	once sync.Once
	repo *repository
)

type IRepository interface {
	Register(form Form) (*entity.Account, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Register(form Form) (*entity.Account, error) {
	var currency entity.Currency
	if err := r.db.Where("symbol=?", form.AccountSymbol).First(&currency).Error; err != nil {
		return nil, utils.NoCurrencyError
	}

	user := entity.User{
		Name:    form.UserName,
		Country: form.UserCountry,
		City:    form.UserCity,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	account := entity.Account{
		UserId:     user.ID,
		CurrencyId: currency.ID,
		Balance:    decimal.NewFromInt32(0),
		DateTime:   time.Now(),
	}
	if err := r.db.Create(&account).Scan(&account).Error; err != nil {
		return nil, err
	}
	account.User = &user
	account.Currency = &currency
	return &account, nil
}

func Init(db *gorm.DB) IRepository {
	once.Do(func() {
		repo = &repository{db: db}
	})
	return repo
}

func GetRepo() IRepository {
	return repo
}
