package entity

import "github.com/google/uuid"

type User struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Name    string    `gorm:"size:256"`
	Country string    `gorm:"size:245"` // TODO: make countries list with foreign key
	City    string    `gorm:"size:256"`
}
