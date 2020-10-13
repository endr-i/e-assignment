package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Rate struct {
	ID         uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Value      decimal.Decimal `gorm:"default:1"`
	CurrencyId uuid.UUID
	Currency   *Currency `gorm:"foreignKey:CurrencyId"`
	Since      time.Time
	Till       time.Time
}
