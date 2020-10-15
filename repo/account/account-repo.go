package accountRepository

import (
	"assignment/entity"
	"assignment/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"sync"
)

var (
	repo *repository
	once sync.Once
)

type IRepository interface {
	Create(entity.User) (*entity.Account, error)
	GetById(string) (*entity.Account, error)
	GetByUuid(uuid2 uuid.UUID) (*entity.Account, error)
}

type repository struct {
	db *gorm.DB
}

func InitRepo(db *gorm.DB) IRepository {
	once.Do(func() {
		repo = &repository{db: db}
	})
	return repo
}

func GetRepo() IRepository {
	return repo
}

func (r *repository) Create(user entity.User) (*entity.Account, error) {
	if user.ID == uuid.Nil {
		return nil, utils.NoUserError
	}
	var account entity.Account
	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) GetByUuid(id uuid.UUID) (*entity.Account, error) {
	if id == uuid.Nil {
		return nil, utils.InvalidUuidError
	}
	var account entity.Account
	if err := r.db.Joins("Currency").Joins("User").First(&account, id).Error; err != nil {
		return nil, utils.NoAccountError
	}
	return &account, nil
}

func (r *repository) GetById(id string) (*entity.Account, error) {
	accountId, err := uuid.Parse(id)
	if err != nil {
		return nil, utils.InvalidDateFormatError
	}
	return r.GetByUuid(accountId)
}

func (r *repository) UpdateBalance(id uuid.UUID, balance decimal.Decimal) error {
	return r.db.Model(&entity.Account{}).Update("balance", balance).Error
}
