package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"player-service/model"
	"player-service/repository"

	"github.com/go-redis/redis/v8"
)

type PlayerService interface {
	CreatePlayer(player *model.Player) error
	GetPlayer(id uint) (*model.Player, error)
	UpdatePlayer(player *model.Player) error
	GetAllPlayers() ([]model.Player, error)
	UpdatePlayerPoints(id uint, delta int) error
}

type playerService struct {
	repo  repository.PlayerRepository
	redis *redis.Client
	ctx   context.Context
}

func NewPlayerService(repo repository.PlayerRepository, redisClient *redis.Client) PlayerService {
	return &playerService{
		repo:  repo,
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *playerService) CreatePlayer(player *model.Player) error {
	player.Power = 0
	return s.repo.Create(player)
}

func (s *playerService) GetPlayer(id uint) (*model.Player, error) {
	cacheKey := fmt.Sprintf("player:%d", id)

	// Try cache first
	cached, err := s.redis.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var player model.Player
		if json.Unmarshal([]byte(cached), &player) == nil {
			return &player, nil
		}
	}

	// Get from database
	player, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	data, _ := json.Marshal(player)
	s.redis.Set(s.ctx, cacheKey, data, 5*time.Minute)

	return player, nil
}

func (s *playerService) UpdatePlayer(player *model.Player) error {
	err := s.repo.Update(player)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("player:%d", player.ID)
	s.redis.Del(s.ctx, cacheKey)

	return nil
}

func (s *playerService) GetAllPlayers() ([]model.Player, error) {
	return s.repo.FindAll()
}

func (s *playerService) UpdatePlayerPoints(id uint, delta int) error {
	player, err := s.GetPlayer(id)
	if err != nil {
		return err
	}

	player.Power += delta
	return s.UpdatePlayer(player)
}
