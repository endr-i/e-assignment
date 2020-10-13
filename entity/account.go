package entity

import (
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
