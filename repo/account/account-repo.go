package accountRepository

import (
	"assignment/entity"
	repo2 "assignment/repo"
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
	GetById(uuid.UUID) (*entity.Account, error)
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
		return nil, repo2.NoUserError
	}
	var account entity.Account
	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) GetById(id uuid.UUID) (*entity.Account, error) {
	if id == uuid.Nil {
		return nil, repo2.NoAccountError
	}
	var account entity.Account
	if err := r.db.Joins("Currency").First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *repository) UpdateBalance(id uuid.UUID, balance decimal.Decimal) error {
	return r.db.Model(&entity.Account{}).Update("balance", balance).Error
}
