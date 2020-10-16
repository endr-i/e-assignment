package pg

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	once sync.Once
	db   *gorm.DB
)

type Config struct {
	Dsn string `default:"user=postgres password=postgresPass host=localhost dbname=e_assignment port=5432 sslmode=disable"`
}

func InitDB(config Config) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(config.Dsn), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			log.Fatal(err)
		}
		db.AutoMigrate(&Migration{})
		version := GetVersion(db)
		migrators := GetMigrators()
		for i := int(version); i < len(migrators); i++ {
			if err := migrators[i].Up(); err != nil {
				break
			}
			version++
		}
		SetVersion(db, version)
	})
	return db
}
