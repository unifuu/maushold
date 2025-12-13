package service

import (
	"errors"
	"math/rand"
	"time"

	"maushold/battle-service/model"
	"maushold/battle-service/repository"

	"github.com/go-redis/redis/v8"
)

type BattleService interface {
	CreateBattle(player1ID, player2ID, pokemon1ID, pokemon2ID uint) (*model.Battle, error)
	GetBattle(id uint) (*model.Battle, error)
	GetPlayerBattles(playerID uint) ([]model.Battle, error)
	GetRecentBattles() ([]model.Battle, error)
}

type battleService struct {
	repo         repository.BattleRepository
	playerClient *PlayerClient
	battleEngine *BattleEngine
	redis        *redis.Client
}

func NewBattleService(
	repo repository.BattleRepository,
	playerClient *PlayerClient,
	battleEngine *BattleEngine,
	redisClient *redis.Client,
) BattleService {
	return &battleService{
		repo:         repo,
		playerClient: playerClient,
		battleEngine: battleEngine,
		redis:        redisClient,
	}
}

func (s *battleService) CreateBattle(player1ID, player2ID, pokemon1ID, pokemon2ID uint) (*model.Battle, error) {
	pokemon1, err := s.playerClient.GetPlayerPokemon(player1ID, pokemon1ID)
	if err != nil {
		return nil, errors.New("pokemon 1 not found")
	}

	pokemon2, err := s.playerClient.GetPlayerPokemon(player2ID, pokemon2ID)
	if err != nil {
		return nil, errors.New("pokemon 2 not found")
	}

	battle := &model.Battle{
		Player1ID:  player1ID,
		Player2ID:  player2ID,
		Pokemon1ID: pokemon1ID,
		Pokemon2ID: pokemon2ID,
		Status:     "in_progress",
	}

	if err := s.repo.Create(battle); err != nil {
		return nil, err
	}

	winner, battleLog := s.battleEngine.SimulateBattle(pokemon1, pokemon2)

	if winner == 1 {
		battle.WinnerID = player1ID
	} else {
		battle.WinnerID = player2ID
	}

	battle.PointsWon = 50 + rand.Intn(50)
	battle.PointsLost = 20 + rand.Intn(30)
	battle.BattleLog = battleLog
	battle.Status = "completed"
	now := time.Now()
	battle.CompletedAt = &now

	s.repo.Update(battle)

	return battle, nil
}

func (s *battleService) GetBattle(id uint) (*model.Battle, error) {
	return s.repo.FindByID(id)
}

func (s *battleService) GetPlayerBattles(playerID uint) ([]model.Battle, error) {
	return s.repo.FindByPlayerID(playerID)
}

func (s *battleService) GetRecentBattles() ([]model.Battle, error) {
	return s.repo.FindRecent(50)
}
