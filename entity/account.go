package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	UserId     uuid.UUID
	User       *User `gorm:"foreignKey:UserId"`
	CurrencyId uuid.UUID
	Currency   *Currency `gorm:"foreignKey:CurrencyId"`
	Balance    decimal.Decimal
}

func (account *Account) MarshalJSON() ([]byte, error) {
	k := decimal.NewFromInt32(100)
	account.Balance = account.Balance.Mul(k).Floor().Div(k)
	return json.Marshal(struct {
		ID       uuid.UUID
		User     *User
		Currency *Currency
		Balance  decimal.Decimal
	}{
		ID:       account.ID,
		User:     account.User,
		Currency: account.Currency,
		Balance:  account.Balance,
	})
}
