package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"maushold/pokemon-service/model"
	"maushold/pokemon-service/repository"

	"github.com/go-redis/redis/v8"
)

type PokemonService interface {
	CreatePokemon(pokemon *model.Pokemon) error
	GetPokemon(id int) (*model.Pokemon, error)
	GetAllPokemon() ([]model.Pokemon, error)
	GetRandomPokemon() (*model.Pokemon, error)
}

type pokemonService struct {
	repo  repository.PokemonRepository
	redis *redis.Client
	ctx   context.Context
}

func NewPokemonService(repo repository.PokemonRepository, redisClient *redis.Client) PokemonService {
	return &pokemonService{
		repo:  repo,
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *pokemonService) CreatePokemon(pokemon *model.Pokemon) error {
	err := s.repo.Create(pokemon)
	if err != nil {
		return err
	}

	s.redis.Del(s.ctx, "pokemon:all")
	return nil
}

func (s *pokemonService) GetPokemon(id int) (*model.Pokemon, error) {
	cacheKey := fmt.Sprintf("pokemon:%d", id)

	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var pokemon model.Pokemon
		if json.Unmarshal([]byte(cached), &pokemon) == nil {
			return &pokemon, nil
		}
	}

	pokemon, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(pokemon)
	s.redis.Set(s.ctx, cacheKey, data, 10*time.Minute)

	return pokemon, nil
}

func (s *pokemonService) GetAllPokemon() ([]model.Pokemon, error) {
	cached, err := s.redis.Get(s.ctx, "pokemon:all").Result()
	if err == nil {
		var pokemon []model.Pokemon
		if json.Unmarshal([]byte(cached), &pokemon) == nil {
			return pokemon, nil
		}
	}

	pokemon, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(pokemon)
	s.redis.Set(s.ctx, "pokemon:all", data, 10*time.Minute)

	return pokemon, nil
}

func (s *pokemonService) GetRandomPokemon() (*model.Pokemon, error) {
	return s.repo.GetRandom()
}
