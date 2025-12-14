package service

import (
	"maushold/player-service/model"
	"maushold/player-service/repository"

	"github.com/go-redis/redis/v8"
)

type PlayerMonsterService interface {
	AddMonsterToPlayer(monster *model.PlayerMonster) error
	GetPlayerMonster(playerID uint) ([]model.PlayerMonster, error)
}

type playerMonsterService struct {
	repo  repository.PlayerMonsterRepository
	redis *redis.Client
}

func NewPlayerMonsterService(repo repository.PlayerMonsterRepository, redisClient *redis.Client) PlayerMonsterService {
	return &playerMonsterService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *playerMonsterService) AddMonsterToPlayer(monster *model.PlayerMonster) error {
	return s.repo.Create(monster)
}

func (s *playerMonsterService) GetPlayerMonster(playerID uint) ([]model.PlayerMonster, error) {
	return s.repo.FindByPlayerID(playerID)
}
