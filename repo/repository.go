package repo

import (
	"assignment/repo/register"
	"gorm.io/gorm"
	"sync"
)

var once sync.Once

func Init(db *gorm.DB) {
	once.Do(func() {
		register.Init(db)
	})
}
