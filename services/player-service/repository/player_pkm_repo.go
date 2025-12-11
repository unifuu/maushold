package repository

import (
	"player-service/model"

	"gorm.io/gorm"
)

type PlayerPokemonRepository interface {
	Create(pokemon *model.PlayerPokemon) error
	FindByPlayerID(playerID uint) ([]model.PlayerPokemon, error)
	FindByID(id uint) (*model.PlayerPokemon, error)
}

type playerPokemonRepository struct {
	db *gorm.DB
}

func NewPlayerPokemonRepository(db *gorm.DB) PlayerPokemonRepository {
	return &playerPokemonRepository{db: db}
}

func (r *playerPokemonRepository) Create(pokemon *model.PlayerPokemon) error {
	return r.db.Create(pokemon).Error
}

func (r *playerPokemonRepository) FindByPlayerID(playerID uint) ([]model.PlayerPokemon, error) {
	var pokemon []model.PlayerPokemon
	err := r.db.Where("player_id = ?", playerID).Find(&pokemon).Error
	return pokemon, err
}

func (r *playerPokemonRepository) FindByID(id uint) (*model.PlayerPokemon, error) {
	var pokemon model.PlayerPokemon
	err := r.db.First(&pokemon, id).Error
	return &pokemon, err
}
