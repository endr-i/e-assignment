package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	OperationId uuid.UUID
	Operation   *Operation `gorm:"foreignKey:OperationId"`
	Value       decimal.Decimal
	AccountId   uuid.UUID
	Account     *Account `gorm:"foreignKey:AccountId"`
	DateTime    time.Time
	CurrencyId  uuid.UUID
	Currency    *Currency `gorm:"foreignKey:CurrencyId"`
	RateValue   decimal.Decimal
}

func (t Transaction) DiffValue() decimal.Decimal {
	return t.Value.Mul(t.RateValue)
}
