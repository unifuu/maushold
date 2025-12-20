package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"maushold/player-service/model"
	"maushold/player-service/repository"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type PlayerService interface {
	CreatePlayer(player *model.Player) error
	GetPlayer(id uint) (*model.Player, error)
	UpdatePlayer(player *model.Player) error
	DeletePlayer(id uint) error
	GetAllPlayers() ([]model.Player, error)
	UpdatePlayerPoints(id uint, delta int) error
	AuthenticatePlayer(username, password string) (*model.Player, error)
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
	player.Points = 0

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(player.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	player.Password = string(hashedPassword)

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

func (s *playerService) DeletePlayer(id uint) error {
	player, err := s.GetPlayer(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(player)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("player:%d", id)
	s.redis.Del(s.ctx, cacheKey)

	return nil
}

func (s *playerService) UpdatePlayerPoints(id uint, delta int) error {
	player, err := s.GetPlayer(id)
	if err != nil {
		return err
	}

	player.Points += delta
	return s.UpdatePlayer(player)
}

func (s *playerService) AuthenticatePlayer(username, password string) (*model.Player, error) {
	// Find player by username
	players, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var player *model.Player
	for i := range players {
		if players[i].Username == username {
			player = &players[i]
			break
		}
	}

	if player == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password with hashed password
	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return player, nil
}
