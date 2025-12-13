package service

import (
	"log"
	"time"

	"maushold/ranking-service/model"
	"maushold/ranking-service/repository"

	"gorm.io/gorm"
)

type RankingService interface {
	UpdatePlayerRanking(playerID uint, pointsDelta int, isWin bool) error
	GetPlayerRanking(playerID uint) (*model.PlayerRanking, error)
	GetLeaderboard(limit int) ([]model.LeaderboardEntry, error)
	SyncRankings() error
	StartPeriodicSync()
}

type rankingService struct {
	repo               repository.RankingRepository
	playerClient       *PlayerClient
	leaderboardService *LeaderboardService
}

func NewRankingService(
	repo repository.RankingRepository,
	playerClient *PlayerClient,
	leaderboardService *LeaderboardService,
) RankingService {
	return &rankingService{
		repo:               repo,
		playerClient:       playerClient,
		leaderboardService: leaderboardService,
	}
}

func (s *rankingService) UpdatePlayerRanking(playerID uint, pointsDelta int, isWin bool) error {
	ranking, err := s.repo.FindByPlayerID(playerID)

	if err == gorm.ErrRecordNotFound {
		player, err := s.playerClient.GetPlayer(playerID)
		if err != nil {
			log.Printf("Player %d not found: %v", playerID, err)
			return err
		}

		ranking = &model.PlayerRanking{
			PlayerID:     playerID,
			Username:     player.Username,
			TotalPoints:  player.Points + pointsDelta,
			TotalBattles: 1,
		}

		if isWin {
			ranking.Wins = 1
		} else {
			ranking.Losses = 1
		}

		err = s.repo.Create(ranking)
	} else if err != nil {
		return err
	} else {
		ranking.TotalPoints += pointsDelta
		ranking.TotalBattles++

		if isWin {
			ranking.Wins++
		} else {
			ranking.Losses++
		}

		err = s.repo.Update(ranking)
	}

	if ranking.TotalBattles > 0 {
		ranking.WinRate = float64(ranking.Wins) / float64(ranking.TotalBattles) * 100
		s.repo.Update(ranking)
	}

	// Update Redis leaderboard
	s.leaderboardService.UpdatePlayerScore(playerID, ranking.TotalPoints)

	log.Printf("Updated ranking for player %d: Points=%d, W/L=%d/%d",
		playerID, ranking.TotalPoints, ranking.Wins, ranking.Losses)

	return err
}

func (s *rankingService) GetPlayerRanking(playerID uint) (*model.PlayerRanking, error) {
	ranking, err := s.repo.FindByPlayerID(playerID)
	if err != nil {
		return nil, err
	}

	// Calculate rank
	allRankings, _ := s.repo.FindAll()
	rank := 1
	for _, r := range allRankings {
		if r.TotalPoints > ranking.TotalPoints {
			rank++
		}
	}
	ranking.Rank = rank

	return ranking, nil
}

func (s *rankingService) GetLeaderboard(limit int) ([]model.LeaderboardEntry, error) {
	// Try Redis first
	leaderboard, err := s.leaderboardService.GetTopPlayers(limit)
	if err == nil && len(leaderboard) > 0 {
		// Enrich with database data
		result := make([]model.LeaderboardEntry, 0, len(leaderboard))
		for i, entry := range leaderboard {
			ranking, err := s.repo.FindByPlayerID(entry.PlayerID)
			if err != nil {
				continue
			}

			result = append(result, model.LeaderboardEntry{
				PlayerID:    entry.PlayerID,
				Username:    ranking.Username,
				TotalPoints: entry.TotalPoints,
				Wins:        ranking.Wins,
				Losses:      ranking.Losses,
				WinRate:     ranking.WinRate,
				Rank:        i + 1,
			})
		}
		return result, nil
	}

	// Fallback to database
	rankings, err := s.repo.FindTopN(limit)
	if err != nil {
		return nil, err
	}

	result := make([]model.LeaderboardEntry, len(rankings))
	for i, r := range rankings {
		result[i] = model.LeaderboardEntry{
			PlayerID:    r.PlayerID,
			Username:    r.Username,
			TotalPoints: r.TotalPoints,
			Wins:        r.Wins,
			Losses:      r.Losses,
			WinRate:     r.WinRate,
			Rank:        i + 1,
		}
	}

	return result, nil
}

func (s *rankingService) SyncRankings() error {
	rankings, err := s.repo.FindAll()
	if err != nil {
		return err
	}

	s.leaderboardService.ClearLeaderboard()

	for _, r := range rankings {
		s.leaderboardService.UpdatePlayerScore(r.PlayerID, r.TotalPoints)
	}

	log.Printf("Ranking sync completed. Synced %d players", len(rankings))
	return nil
}

func (s *rankingService) StartPeriodicSync() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.SyncRankings()
	}
}
