package service

import (
	"player-service/model"
	"player-service/repository"

	"github.com/go-redis/redis/v8"
)

type PlayerPokemonService interface {
	AddPokemonToPlayer(pokemon *model.PlayerPokemon) error
	GetPlayerPokemon(playerID uint) ([]model.PlayerPokemon, error)
}

type playerPokemonService struct {
	repo  repository.PlayerPokemonRepository
	redis *redis.Client
}

func NewPlayerPokemonService(repo repository.PlayerPokemonRepository, redisClient *redis.Client) PlayerPokemonService {
	return &playerPokemonService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *playerPokemonService) AddPokemonToPlayer(pokemon *model.PlayerPokemon) error {
	return s.repo.Create(pokemon)
}

func (s *playerPokemonService) GetPlayerPokemon(playerID uint) ([]model.PlayerPokemon, error) {
	return s.repo.FindByPlayerID(playerID)
}
