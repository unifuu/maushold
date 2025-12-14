package repository

import (
	"maushold/monster-service/model"

	"gorm.io/gorm"
)

type MonsterRepository interface {
	Create(monster *model.Monster) error
	FindByID(id int) (*model.Monster, error)
	FindAll() ([]model.Monster, error)
	GetRandom() (*model.Monster, error)
}

type monsterRepository struct {
	db *gorm.DB
}

func NewMonsterRepository(db *gorm.DB) MonsterRepository {
	return &monsterRepository{db: db}
}

func (r *monsterRepository) Create(monster *model.Monster) error {
	return r.db.Create(monster).Error
}

func (r *monsterRepository) FindByID(id int) (*model.Monster, error) {
	var monster model.Monster
	err := r.db.First(&monster, id).Error
	return &monster, err
}

func (r *monsterRepository) FindAll() ([]model.Monster, error) {
	var monster []model.Monster
	err := r.db.Find(&monster).Error
	return monster, err
}

func (r *monsterRepository) GetRandom() (*model.Monster, error) {
	var monster model.Monster
	err := r.db.Order("RANDOM()").First(&monster).Error
	return &monster, err
}
