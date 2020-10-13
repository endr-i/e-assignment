package entity

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type Operation struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Type     int
	Details  datatypes.JSON
	DateTime time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
