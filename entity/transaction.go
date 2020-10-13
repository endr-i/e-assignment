package entity

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	OperationId uuid.UUID
	Operation   *Operation `gorm:"foreignKey:OperationId"`
	Value       decimal.Decimal
}
