package repository

import (
	"maushold/pokemon-service/model"

	"gorm.io/gorm"
)

type PokemonRepository interface {
	Create(pokemon *model.Pokemon) error
	FindByID(id int) (*model.Pokemon, error)
	FindAll() ([]model.Pokemon, error)
	GetRandom() (*model.Pokemon, error)
}

type pokemonRepository struct {
	db *gorm.DB
}

func NewPokemonRepository(db *gorm.DB) PokemonRepository {
	return &pokemonRepository{db: db}
}

func (r *pokemonRepository) Create(pokemon *model.Pokemon) error {
	return r.db.Create(pokemon).Error
}

func (r *pokemonRepository) FindByID(id int) (*model.Pokemon, error) {
	var pokemon model.Pokemon
	err := r.db.First(&pokemon, id).Error
	return &pokemon, err
}

func (r *pokemonRepository) FindAll() ([]model.Pokemon, error) {
	var pokemon []model.Pokemon
	err := r.db.Find(&pokemon).Error
	return pokemon, err
}

func (r *pokemonRepository) GetRandom() (*model.Pokemon, error) {
	var pokemon model.Pokemon
	err := r.db.Order("RANDOM()").First(&pokemon).Error
	return &pokemon, err
}
