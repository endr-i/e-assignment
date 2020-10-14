package entity

import (
	"github.com/google/uuid"
)

type Currency struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Name   string    `gorm:"default:"`
	Symbol string    `gorm:"unique"`
}
