package pg

import (
	"assignment/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Migration struct {
	Version int32 `gorm:"default:0"`
}

type Migrator interface {
	Up() error
	Down() error
}

type initMigration struct {
}

func (initMigration) Up() error {
	return db.Transaction(func(tx *gorm.DB) error {
		// TODO: make correct migrations
		err := tx.AutoMigrate(&entity.Currency{}, &entity.User{})
		if err != nil {
			return err
		}
		return tx.AutoMigrate(&entity.Rate{}, &entity.Account{}, &entity.Transaction{})
	})
}

func (initMigration) Down() error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Migrator().DropTable(
			&entity.Currency{},
			&entity.User{},
			&entity.Rate{},
			&entity.Account{},
			&entity.Transaction{},
		)
	})
}

type initDataMigration struct {
}

func (initDataMigration) Up() error {
	return db.Transaction(func(tx *gorm.DB) error {
		usdUuid := uuid.New()
		tx.Create(&entity.Currency{
			ID:     usdUuid,
			Name:   "US Dollar",
			Symbol: "USD",
		})
		eurUuid := uuid.New()
		tx.Create(&entity.Currency{
			ID:     eurUuid,
			Name:   "Euro",
			Symbol: "EUR",
		})
		rubUuid := uuid.New()
		tx.Create(&entity.Currency{
			ID:     rubUuid,
			Name:   "Russian ruble",
			Symbol: "RUB",
		})
		now := time.Now()
		tx.Create(&entity.Rate{
			CurrencyId: usdUuid,
			Value:      decimal.NewFromFloat(1),
			DateTime:   now,
		})
		tx.Create(&entity.Rate{
			CurrencyId: rubUuid,
			Value:      decimal.NewFromFloat(0.013),
			DateTime:   now,
		})
		tx.Create(&entity.Rate{
			CurrencyId: eurUuid,
			Value:      decimal.NewFromFloat(1.17),
			DateTime:   now,
		})
		return nil
	})
}

func (initDataMigration) Down() error {
	return db.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&entity.Rate{}).Where("1 = 1")
		tx.Delete(&entity.Currency{}).Where("1 = 1")
		return nil
	})
}

func GetMigrators() []Migrator {
	return []Migrator{
		initMigration{},
		initDataMigration{},
	}
}

func GetVersion(db *gorm.DB) int32 {
	var migration Migration
	db.FirstOrCreate(&migration)
	return migration.Version
}

func SetVersion(db *gorm.DB, version int32) error {
	return db.Exec(`Update "migrations" SET "version"=?`, version).Error
}
