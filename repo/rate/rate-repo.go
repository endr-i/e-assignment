package rateRepository

import (
	"assignment/entity"
	"assignment/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	repo         *repository
	once         sync.Once
	mainCurrency *entity.Currency
)

type IRepository interface {
	GetMainCurrency() *entity.Currency
	GetCurrencyRate(uuid.UUID, time.Time) (*entity.Rate, error)
	//GetConvertRate(uuid.UUID, uuid.UUID, time.Time) (*ConvertRate, error)
	UploadRates(UploadRatesForm) ([]entity.Rate, error)
	Create(CreateForm, time.Time) (*entity.Rate, error)
}

type repository struct {
	db *gorm.DB
}

func InitRepo(db *gorm.DB) IRepository {
	once.Do(func() {
		repo = &repository{db: db}
		var usd entity.Currency
		db.First(&usd, "symbol=?", "USD")
		mainCurrency = &usd
	})
	return repo
}

func GetRepo() IRepository {
	return repo
}

func (r *repository) GetMainCurrency() *entity.Currency {
	return mainCurrency
}

func (r *repository) GetCurrencyRate(currencyId uuid.UUID, rateTime time.Time) (*entity.Rate, error) {
	if currencyId == uuid.Nil {
		return nil, utils.InvalidUuidError
	}
	var rate entity.Rate
	err := r.db.Where("currency_id=?", currencyId, rateTime, rateTime).
		Order("date_time DESC").First(&rate).Error
	if err != nil {
		return nil, utils.NoRateError
	}
	return &rate, nil
}

//func (r *repository) GetConvertRate(fromId uuid.UUID, toId uuid.UUID, rateTime time.Time) (*ConvertRate, error) {
//	rate1, err := r.GetCurrencyRate(fromId, rateTime)
//	if err != nil {
//		return nil, err
//	}
//	rate2, err := r.GetCurrencyRate(toId, rateTime)
//	if err != nil {
//		return nil, err
//	}
//	return &ConvertRate{
//		RateFrom:  *rate1,
//		RateTo:    *rate2,
//		RateValue: utils.ConvertRate(rate1.Value, rate2.Value),
//	}, nil
//}

func (r *repository) Create(form CreateForm, rateTime time.Time) (*entity.Rate, error) {
	var currency entity.Currency
	if err := r.db.First(&currency, "symbol=?", form.Symbol).Error; err != nil {
		return nil, utils.NoCurrencyError
	}
	rate := entity.Rate{
		Value:      decimal.NewFromFloat(form.Value),
		CurrencyId: currency.ID,
		Currency:   &currency,
		DateTime:   rateTime,
	}
	if err := r.db.Create(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *repository) UploadRates(form UploadRatesForm) ([]entity.Rate, error) {
	if form.Rates == nil {
		return nil, utils.NoRatesToUploadError
	}
	now := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(form.Rates))
	rates := make([]*entity.Rate, len(form.Rates))
	for i, formRate := range form.Rates {
		go func(formRate CreateForm, i int) {
			defer wg.Done()
			rate, err := r.Create(formRate, now)
			if err == nil {
				rates[i] = rate
			}
		}(formRate, i)
	}
	wg.Wait()

	result := make([]entity.Rate, 0, len(rates))
	for _, rate := range rates {
		if rate != nil {
			result = append(result, *rate)
		}
	}
	return result, nil
}
