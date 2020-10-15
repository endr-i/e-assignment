package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	OperationId       uuid.UUID
	Operation         *Operation `gorm:"foreignKey:OperationId"`
	Value             decimal.Decimal
	AccountId         uuid.UUID
	Account           *Account `gorm:"foreignKey:AccountId"`
	DateTime          time.Time
	CurrencyId        uuid.UUID
	Currency          *Currency `gorm:"foreignKey:CurrencyId"`
	CurrencyRateValue decimal.Decimal
	AccountRateValue  decimal.Decimal
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	value, _ := t.Value.Float64()
	return json.Marshal(struct {
		ID       uuid.UUID
		Value    float64
		Account  *Account
		DateTime time.Time
		Currency *Currency
	}{
		ID:       t.ID,
		Value:    value,
		Account:  t.Account,
		DateTime: t.DateTime,
		Currency: t.Currency,
	})
}

func (t Transaction) RatedValue() decimal.Decimal {
	return t.UsdValue().Div(t.AccountRateValue)
}

func (t Transaction) UsdValue() decimal.Decimal {
	return t.Value.Mul(t.CurrencyRateValue)
}
