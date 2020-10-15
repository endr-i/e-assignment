package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Account struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	UserId     uuid.UUID
	User       *User `gorm:"foreignKey:UserId"`
	CurrencyId uuid.UUID
	Currency   *Currency `gorm:"foreignKey:CurrencyId"`
	Balance    decimal.Decimal
	DateTime   time.Time
}

func (account *Account) MarshalJSON() ([]byte, error) {
	balance, _ := account.Balance.Round(2).Float64()
	return json.Marshal(struct {
		ID       uuid.UUID
		User     *User
		Currency *Currency
		Balance  float64
	}{
		ID:       account.ID,
		User:     account.User,
		Currency: account.Currency,
		Balance:  balance,
	})
}
