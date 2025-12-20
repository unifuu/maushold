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
	UpdatePlayerCombatPower(playerID uint, combatPower int64) error
	GetPlayerRanking(playerID uint) (*model.PlayerRanking, error)
	GetLeaderboard(limit int) (*model.LeaderboardResponse, error)
	GetPlayerRankWithContext(playerID uint, contextSize int) (*model.PlayerRankContext, error)
	SyncRankings() error
	StartPeriodicSync()
	RefreshMaterializedView() error
	SyncFromPlayerService() error
}

type rankingService struct {
	repo               repository.RankingRepository
	playerClient       *PlayerClient
	battleClient       *BattleClient
	leaderboardService *LeaderboardService
}

func NewRankingService(
	repo repository.RankingRepository,
	playerClient *PlayerClient,
	battleClient *BattleClient,
	leaderboardService *LeaderboardService,
) RankingService {
	return &rankingService{
		repo:               repo,
		playerClient:       playerClient,
		battleClient:       battleClient,
		leaderboardService: leaderboardService,
	}
}

// UpdatePlayerRanking updates player stats after a battle
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
			CombatPower:  int64(player.Points+pointsDelta) * 100, // Combat power = points * 100
			TotalBattles: 1,
			LastBattleAt: time.Now(),
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
		ranking.CombatPower = int64(ranking.TotalPoints) * 100 // Update combat power
		ranking.TotalBattles++
		ranking.LastBattleAt = time.Now()

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

	// Update Redis leaderboard with combat power (threshold-based)
	s.leaderboardService.UpdatePlayerScore(playerID, ranking.CombatPower)

	// Cache player details in Redis
	s.leaderboardService.CachePlayerDetails(ranking)

	log.Printf("Updated ranking for player %d: CombatPower=%d, Points=%d, W/L=%d/%d",
		playerID, ranking.CombatPower, ranking.TotalPoints, ranking.Wins, ranking.Losses)

	return err
}

// UpdatePlayerCombatPower directly updates a player's combat power
func (s *rankingService) UpdatePlayerCombatPower(playerID uint, combatPower int64) error {
	err := s.repo.UpdateCombatPower(playerID, combatPower)
	if err != nil {
		return err
	}

	// Update Redis leaderboard
	err = s.leaderboardService.UpdatePlayerScore(playerID, combatPower)
	if err != nil {
		log.Printf("Failed to update Redis for player %d: %v", playerID, err)
	}

	// Refresh cached player details
	ranking, err := s.repo.FindByPlayerID(playerID)
	if err == nil {
		s.leaderboardService.CachePlayerDetails(ranking)
	}

	return nil
}

// GetPlayerRanking returns a player's ranking information
func (s *rankingService) GetPlayerRanking(playerID uint) (*model.PlayerRanking, error) {
	// Try Redis cache first (fastest for detail lookup)
	cachedPlayer, err := s.leaderboardService.GetPlayerDetails(playerID)
	if err == nil && cachedPlayer != nil {
		// Also try to get rank from Redis
		rank, err := s.leaderboardService.GetPlayerRank(playerID)
		if err == nil && rank > 0 {
			cachedPlayer.Rank = int(rank)
			return cachedPlayer, nil
		}
	}

	// Fallback to DB
	ranking, err := s.repo.FindByPlayerID(playerID)
	if err != nil {
		return nil, err
	}

	// Try to get rank from Redis
	redisRank, err := s.leaderboardService.GetPlayerRank(playerID)
	if err == nil && redisRank > 0 {
		ranking.Rank = int(redisRank)
		return ranking, nil
	}

	// Try to get rank from materialized view
	rank, err := s.repo.GetPlayerRankFromMaterializedView(playerID)
	if err == nil && rank > 0 {
		ranking.Rank = rank
		return ranking, nil
	}

	// Last resort: calculate from database
	allRankings, _ := s.repo.FindAll()
	rank = 1
	for _, r := range allRankings {
		if r.CombatPower > ranking.CombatPower {
			rank++
		}
	}
	ranking.Rank = rank

	return ranking, nil
}

// GetLeaderboard returns the top N players with metadata
func (s *rankingService) GetLeaderboard(limit int) (*model.LeaderboardResponse, error) {
	// Try Redis cache first (highest performance for high traffic)
	redisEntries, err := s.leaderboardService.GetTopPlayers(limit)
	if err == nil && len(redisEntries) > 0 {
		// Enrich with player details from cache
		enrichedEntries := s.enrichLeaderboardEntries(redisEntries)
		metadata, _ := s.getMetadata(true)
		return &model.LeaderboardResponse{
			Leaderboard: enrichedEntries,
			Metadata:    *metadata,
		}, nil
	}

	// Fallback 1: Try materialized view (fast persistent storage)
	entries, err := s.repo.FindTopNFromMaterializedView(limit)
	if err == nil && len(entries) > 0 {
		metadata, _ := s.getMetadata(true)
		return &model.LeaderboardResponse{
			Leaderboard: entries,
			Metadata:    *metadata,
		}, nil
	}

	// Fallback 2: Direct database query (consistent but slower)
	rankings, err := s.repo.FindTopNByCombatPower(limit)
	if err != nil {
		return nil, err
	}

	entries = make([]model.LeaderboardEntry, len(rankings))
	for i, r := range rankings {
		entries[i] = model.LeaderboardEntry{
			PlayerID:    r.PlayerID,
			Username:    r.Username,
			CombatPower: r.CombatPower,
			TotalPoints: r.TotalPoints,
			Wins:        r.Wins,
			Losses:      r.Losses,
			WinRate:     r.WinRate,
			Rank:        i + 1,
			UpdatedAt:   r.UpdatedAt,
		}
	}

	metadata, _ := s.getMetadata(false)
	return &model.LeaderboardResponse{
		Leaderboard: entries,
		Metadata:    *metadata,
	}, nil
}

// GetPlayerRankWithContext returns a player's rank with surrounding players
func (s *rankingService) GetPlayerRankWithContext(playerID uint, contextSize int) (*model.PlayerRankContext, error) {
	// Get player's rank
	rank, err := s.leaderboardService.GetPlayerRank(playerID)
	if err != nil || rank == 0 {
		// Player not in top 10K, get from DB
		ranking, err := s.repo.FindByPlayerID(playerID)
		if err != nil {
			return nil, err
		}

		metadata, _ := s.getMetadata(false)
		return &model.PlayerRankContext{
			Player: model.LeaderboardEntry{
				PlayerID:    ranking.PlayerID,
				Username:    ranking.Username,
				CombatPower: ranking.CombatPower,
				TotalPoints: ranking.TotalPoints,
				Wins:        ranking.Wins,
				Losses:      ranking.Losses,
				WinRate:     ranking.WinRate,
				Rank:        0, // Not in top 10K
				UpdatedAt:   ranking.UpdatedAt,
			},
			Neighbors: []model.LeaderboardEntry{},
			Metadata:  *metadata,
		}, nil
	}

	// Get surrounding players from Redis
	neighbors, err := s.leaderboardService.GetPlayersAroundRank(rank, contextSize)
	if err != nil {
		return nil, err
	}

	// Enrich with player details
	enrichedNeighbors := s.enrichLeaderboardEntries(neighbors)

	// Find the player in neighbors
	var playerEntry model.LeaderboardEntry
	for _, entry := range enrichedNeighbors {
		if entry.PlayerID == playerID {
			playerEntry = entry
			break
		}
	}

	metadata, _ := s.getMetadata(true)
	return &model.PlayerRankContext{
		Player:    playerEntry,
		Neighbors: enrichedNeighbors,
		Metadata:  *metadata,
	}, nil
}

// SyncRankings syncs all rankings from DB to Redis (full rebuild)
func (s *rankingService) SyncRankings() error {
	// Acquire distributed lock
	locked, err := s.leaderboardService.AcquireSyncLock()
	if err != nil || !locked {
		log.Printf("Could not acquire sync lock, another instance is syncing")
		return err
	}
	defer s.leaderboardService.ReleaseSyncLock()

	log.Printf("Starting full ranking sync...")

	// Get top 10K players from database
	rankings, err := s.repo.FindTopNByCombatPower(Top10KLimit)
	if err != nil {
		return err
	}

	// Clear Redis leaderboard
	s.leaderboardService.ClearLeaderboard()

	// Batch update Redis
	err = s.leaderboardService.BatchUpdatePlayerScores(rankings)
	if err != nil {
		return err
	}

	// Cache player details
	for _, ranking := range rankings {
		s.leaderboardService.CachePlayerDetails(&ranking)
	}

	// Update metadata
	totalPlayers, _ := s.repo.GetTotalPlayerCount()
	s.leaderboardService.SetMetadata(totalPlayers)

	log.Printf("Ranking sync completed. Synced %d players to Redis", len(rankings))
	return nil
}

// RefreshMaterializedView refreshes the materialized view
func (s *rankingService) RefreshMaterializedView() error {
	log.Printf("Refreshing materialized view...")
	err := s.repo.RefreshMaterializedView()
	if err != nil {
		log.Printf("Failed to refresh materialized view: %v", err)
		return err
	}
	log.Printf("Materialized view refreshed successfully")
	return nil
}

// SyncFromPlayerService pulls all players from player-service and populates rankings
func (s *rankingService) SyncFromPlayerService() error {
	log.Printf("Starting initial data sync from player-service and battle-service...")
	players, err := s.playerClient.GetAllPlayers()
	if err != nil {
		log.Printf("Failed to fetch players from player-service: %v", err)
		return err
	}

	for _, p := range players {
		ranking, err := s.repo.FindByPlayerID(p.ID)
		if err == gorm.ErrRecordNotFound {
			ranking = &model.PlayerRanking{
				PlayerID:    p.ID,
				Username:    p.Username,
				TotalPoints: p.Points,
				CombatPower: int64(p.Points) * 100,
			}
			s.repo.Create(ranking)
		} else if err == nil {
			ranking.TotalPoints = p.Points
			ranking.CombatPower = int64(p.Points) * 100
			s.repo.Update(ranking)
		}
	}

	// Sync battle history for win/loss stats
	battles, err := s.battleClient.GetAllBattles()
	if err != nil {
		log.Printf("Failed to fetch battles from battle-service: %v", err)
		// Don't return, we still synced points
	} else {
		log.Printf("Syncing %d battles from history...", len(battles))
		for _, b := range battles {
			if b.Status != "completed" {
				continue
			}

			// Update winner
			if b.WinnerID != 0 {
				s.UpdatePlayerRanking(b.WinnerID, 0, true)
			}

			// Update loser
			var loserID uint
			if b.WinnerID == b.Player1ID {
				loserID = b.Player2ID
			} else {
				loserID = b.Player1ID
			}

			if loserID != 0 {
				s.UpdatePlayerRanking(loserID, 0, false)
			}
		}
	}

	log.Printf("Initial data sync completed. Processed %d players and %d battles.", len(players), len(battles))
	return nil
}

// StartPeriodicSync starts periodic sync tasks
func (s *rankingService) StartPeriodicSync() {
	// Sync Redis every 5 minutes
	redisSyncTicker := time.NewTicker(5 * time.Minute)
	go func() {
		for range redisSyncTicker.C {
			s.SyncRankings()
		}
	}()

	// Refresh materialized view every 2 minutes
	mvRefreshTicker := time.NewTicker(2 * time.Minute)
	go func() {
		for range mvRefreshTicker.C {
			s.RefreshMaterializedView()
		}
	}()

	log.Printf("Periodic sync tasks started")
}

// Helper: enrich leaderboard entries with player details from cache
func (s *rankingService) enrichLeaderboardEntries(entries []model.LeaderboardEntry) []model.LeaderboardEntry {
	if len(entries) == 0 {
		return entries
	}

	// Extract player IDs
	playerIDs := make([]uint, len(entries))
	for i, entry := range entries {
		playerIDs[i] = entry.PlayerID
	}

	// Batch get player details from cache
	cachedPlayers, err := s.leaderboardService.BatchGetPlayerDetails(playerIDs)
	if err != nil {
		// If cache fails, try database
		dbPlayers, err := s.repo.FindByPlayerIDs(playerIDs)
		if err == nil {
			cachedPlayers = make(map[uint]*model.PlayerRanking)
			for i := range dbPlayers {
				cachedPlayers[dbPlayers[i].PlayerID] = &dbPlayers[i]
			}
		}
	}

	// Enrich entries
	for i := range entries {
		if player, ok := cachedPlayers[entries[i].PlayerID]; ok {
			entries[i].Username = player.Username
			entries[i].TotalPoints = player.TotalPoints
			entries[i].Wins = player.Wins
			entries[i].Losses = player.Losses
			entries[i].WinRate = player.WinRate
			entries[i].UpdatedAt = player.UpdatedAt
		}
	}

	return entries
}

// Helper: get metadata
func (s *rankingService) getMetadata(cacheHit bool) (*model.LeaderboardMetadata, error) {
	metadata, err := s.leaderboardService.GetMetadata()
	if err != nil {
		// Fallback to database
		totalPlayers, _ := s.repo.GetTotalPlayerCount()
		threshold, _ := s.repo.GetTop10KThreshold()
		metadata = &model.LeaderboardMetadata{
			TotalPlayers:    totalPlayers,
			Top10KThreshold: threshold,
			LastUpdated:     time.Now(),
			CacheHit:        false,
		}
	} else {
		metadata.CacheHit = cacheHit
	}
	return metadata, nil
}
