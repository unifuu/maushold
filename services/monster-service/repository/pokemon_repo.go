package repository

import (
	"maushold/monster-service/model"

	"gorm.io/gorm"
)

type PokemonRepository interface {
	Create(monster *model.Pokemon) error
	FindByID(id int) (*model.Pokemon, error)
	FindAll() ([]model.Pokemon, error)
	GetRandom() (*model.Pokemon, error)
}

type monsterRepository struct {
	db *gorm.DB
}

func NewPokemonRepository(db *gorm.DB) PokemonRepository {
	return &monsterRepository{db: db}
}

func (r *monsterRepository) Create(monster *model.Pokemon) error {
	return r.db.Create(monster).Error
}

func (r *monsterRepository) FindByID(id int) (*model.Pokemon, error) {
	var monster model.Pokemon
	err := r.db.First(&monster, id).Error
	return &monster, err
}

func (r *monsterRepository) FindAll() ([]model.Pokemon, error) {
	var monster []model.Pokemon
	err := r.db.Find(&monster).Error
	return monster, err
}

func (r *monsterRepository) GetRandom() (*model.Pokemon, error) {
	var monster model.Pokemon
	err := r.db.Order("RANDOM()").First(&monster).Error
	return &monster, err
}
