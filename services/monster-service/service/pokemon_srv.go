package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"maushold/monster-service/model"
	"maushold/monster-service/repository"

	"github.com/go-redis/redis/v8"
)

type PokemonService interface {
	CreatePokemon(monster *model.Pokemon) error
	GetPokemon(id int) (*model.Pokemon, error)
	GetAllPokemon() ([]model.Pokemon, error)
	GetRandomPokemon() (*model.Pokemon, error)
}

type monsterService struct {
	repo  repository.PokemonRepository
	redis *redis.Client
	ctx   context.Context
}

func NewPokemonService(repo repository.PokemonRepository, redisClient *redis.Client) PokemonService {
	return &monsterService{
		repo:  repo,
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *monsterService) CreatePokemon(monster *model.Pokemon) error {
	err := s.repo.Create(monster)
	if err != nil {
		return err
	}

	s.redis.Del(s.ctx, "monster:all")
	return nil
}

func (s *monsterService) GetPokemon(id int) (*model.Pokemon, error) {
	cacheKey := fmt.Sprintf("monster:%d", id)

	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var monster model.Pokemon
		if json.Unmarshal([]byte(cached), &monster) == nil {
			return &monster, nil
		}
	}

	monster, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(monster)
	s.redis.Set(s.ctx, cacheKey, data, 10*time.Minute)

	return monster, nil
}

func (s *monsterService) GetAllPokemon() ([]model.Pokemon, error) {
	cached, err := s.redis.Get(s.ctx, "monster:all").Result()
	if err == nil {
		var monster []model.Pokemon
		if json.Unmarshal([]byte(cached), &monster) == nil {
			return monster, nil
		}
	}

	monster, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(monster)
	s.redis.Set(s.ctx, "monster:all", data, 10*time.Minute)

	return monster, nil
}

func (s *monsterService) GetRandomPokemon() (*model.Pokemon, error) {
	return s.repo.GetRandom()
}
