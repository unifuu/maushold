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

type MonsterService interface {
	CreateMonster(monster *model.Monster) error
	GetMonster(id int) (*model.Monster, error)
	GetAllMonster() ([]model.Monster, error)
	GetRandomMonster() (*model.Monster, error)
}

type monsterService struct {
	repo  repository.MonsterRepository
	redis *redis.Client
	ctx   context.Context
}

func NewMonsterService(repo repository.MonsterRepository, redisClient *redis.Client) MonsterService {
	return &monsterService{
		repo:  repo,
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *monsterService) CreateMonster(monster *model.Monster) error {
	err := s.repo.Create(monster)
	if err != nil {
		return err
	}

	s.redis.Del(s.ctx, "monster:all")
	return nil
}

func (s *monsterService) GetMonster(id int) (*model.Monster, error) {
	cacheKey := fmt.Sprintf("monster:%d", id)

	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var monster model.Monster
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

func (s *monsterService) GetAllMonster() ([]model.Monster, error) {
	cached, err := s.redis.Get(s.ctx, "monster:all").Result()
	if err == nil {
		var monster []model.Monster
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

func (s *monsterService) GetRandomMonster() (*model.Monster, error) {
	return s.repo.GetRandom()
}
