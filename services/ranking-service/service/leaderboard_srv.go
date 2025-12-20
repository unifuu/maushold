package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"maushold/ranking-service/model"

	"github.com/go-redis/redis/v8"
)

const (
	LeaderboardKey      = "leaderboard:global"
	ThresholdKey        = "leaderboard:threshold"
	TotalPlayersKey     = "leaderboard:total_players"
	LastSyncKey         = "leaderboard:last_sync"
	SyncLockKey         = "leaderboard:sync_lock"
	PlayerDetailsPrefix = "player:"
	Top10KLimit         = 10000
	LockTimeout         = 5 * time.Minute
	PlayerCacheTTL      = 24 * time.Hour
)

type LeaderboardService struct {
	redis      *redis.Client
	ctx        context.Context
	instanceID string
}

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{
		redis:      redisClient,
		ctx:        context.Background(),
		instanceID: fmt.Sprintf("instance-%d", time.Now().UnixNano()),
	}
}

// UpdatePlayerScore updates a player's score with threshold checking
func (s *LeaderboardService) UpdatePlayerScore(playerID uint, combatPower int64) error {
	// If combat power is 0 or less, remove from leaderboard
	if combatPower <= 0 {
		return s.redis.ZRem(s.ctx, LeaderboardKey, fmt.Sprintf("%d", playerID)).Err()
	}

	// Get current threshold
	threshold, err := s.GetThreshold()
	if err != nil {
		threshold = 0 // If no threshold, allow update
	}

	// Only update Redis if player is in or above top 10K threshold
	if combatPower >= threshold {
		// Add to sorted set
		err := s.redis.ZAdd(s.ctx, LeaderboardKey, &redis.Z{
			Score:  float64(combatPower),
			Member: fmt.Sprintf("%d", playerID),
		}).Err()
		if err != nil {
			return err
		}

		// Trim to top 10K
		err = s.redis.ZRemRangeByRank(s.ctx, LeaderboardKey, 0, -Top10KLimit-1).Err()
		if err != nil {
			return err
		}

		// Update threshold
		return s.updateThreshold()
	}

	return nil
}

// UpdatePlayerScoreForce forces update regardless of threshold (for initial sync)
func (s *LeaderboardService) UpdatePlayerScoreForce(playerID uint, combatPower int64) error {
	return s.redis.ZAdd(s.ctx, LeaderboardKey, &redis.Z{
		Score:  float64(combatPower),
		Member: fmt.Sprintf("%d", playerID),
	}).Err()
}

// BatchUpdatePlayerScores updates multiple players efficiently
func (s *LeaderboardService) BatchUpdatePlayerScores(players []model.PlayerRanking) error {
	if len(players) == 0 {
		return nil
	}

	pipe := s.redis.Pipeline()

	for _, player := range players {
		pipe.ZAdd(s.ctx, LeaderboardKey, &redis.Z{
			Score:  float64(player.CombatPower),
			Member: fmt.Sprintf("%d", player.PlayerID),
		})
	}

	_, err := pipe.Exec(s.ctx)
	if err != nil {
		return err
	}

	// Trim to top 10K
	return s.redis.ZRemRangeByRank(s.ctx, LeaderboardKey, 0, -Top10KLimit-1).Err()
}

// GetTopPlayers returns top N players from Redis
func (s *LeaderboardService) GetTopPlayers(limit int) ([]model.LeaderboardEntry, error) {
	result, err := s.redis.ZRevRangeWithScores(s.ctx, LeaderboardKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]model.LeaderboardEntry, 0, len(result))
	for i, z := range result {
		playerID, _ := strconv.ParseUint(z.Member.(string), 10, 64)
		entries = append(entries, model.LeaderboardEntry{
			PlayerID:    uint(playerID),
			CombatPower: int64(z.Score),
			Rank:        i + 1,
		})
	}

	return entries, nil
}

// GetPlayerRank returns the rank of a specific player (1-indexed)
func (s *LeaderboardService) GetPlayerRank(playerID uint) (int64, error) {
	rank, err := s.redis.ZRevRank(s.ctx, LeaderboardKey, fmt.Sprintf("%d", playerID)).Result()
	if err == redis.Nil {
		return 0, nil // Player not in top 10K
	}
	if err != nil {
		return 0, err
	}
	return rank + 1, nil // Convert to 1-indexed
}

// GetPlayerScore returns the combat power of a specific player
func (s *LeaderboardService) GetPlayerScore(playerID uint) (int64, error) {
	score, err := s.redis.ZScore(s.ctx, LeaderboardKey, fmt.Sprintf("%d", playerID)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int64(score), nil
}

// GetPlayersAroundRank returns players around a specific rank (for context)
func (s *LeaderboardService) GetPlayersAroundRank(rank int64, context int) ([]model.LeaderboardEntry, error) {
	start := rank - int64(context) - 1
	if start < 0 {
		start = 0
	}
	end := rank + int64(context) - 1

	result, err := s.redis.ZRevRangeWithScores(s.ctx, LeaderboardKey, start, end).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]model.LeaderboardEntry, 0, len(result))
	for i, z := range result {
		playerID, _ := strconv.ParseUint(z.Member.(string), 10, 64)
		entries = append(entries, model.LeaderboardEntry{
			PlayerID:    uint(playerID),
			CombatPower: int64(z.Score),
			Rank:        int(start) + i + 1,
		})
	}

	return entries, nil
}

// CachePlayerDetails caches player details in Redis hash
func (s *LeaderboardService) CachePlayerDetails(player *model.PlayerRanking) error {
	key := fmt.Sprintf("%s%d", PlayerDetailsPrefix, player.PlayerID)

	data := map[string]interface{}{
		"username":     player.Username,
		"combat_power": player.CombatPower,
		"total_points": player.TotalPoints,
		"wins":         player.Wins,
		"losses":       player.Losses,
		"win_rate":     player.WinRate,
		"updated_at":   player.UpdatedAt.Format(time.RFC3339),
	}

	pipe := s.redis.Pipeline()
	pipe.HSet(s.ctx, key, data)
	pipe.Expire(s.ctx, key, PlayerCacheTTL)
	_, err := pipe.Exec(s.ctx)

	return err
}

// GetPlayerDetails retrieves cached player details
func (s *LeaderboardService) GetPlayerDetails(playerID uint) (*model.PlayerRanking, error) {
	key := fmt.Sprintf("%s%d", PlayerDetailsPrefix, playerID)

	result, err := s.redis.HGetAll(s.ctx, key).Result()
	if err != nil || len(result) == 0 {
		return nil, err
	}

	combatPower, _ := strconv.ParseInt(result["combat_power"], 10, 64)
	totalPoints, _ := strconv.Atoi(result["total_points"])
	wins, _ := strconv.Atoi(result["wins"])
	losses, _ := strconv.Atoi(result["losses"])
	winRate, _ := strconv.ParseFloat(result["win_rate"], 64)
	updatedAt, _ := time.Parse(time.RFC3339, result["updated_at"])

	return &model.PlayerRanking{
		PlayerID:    playerID,
		Username:    result["username"],
		CombatPower: combatPower,
		TotalPoints: totalPoints,
		Wins:        wins,
		Losses:      losses,
		WinRate:     winRate,
		UpdatedAt:   updatedAt,
	}, nil
}

// BatchGetPlayerDetails retrieves multiple player details efficiently
func (s *LeaderboardService) BatchGetPlayerDetails(playerIDs []uint) (map[uint]*model.PlayerRanking, error) {
	if len(playerIDs) == 0 {
		return make(map[uint]*model.PlayerRanking), nil
	}

	pipe := s.redis.Pipeline()
	cmds := make(map[uint]*redis.StringStringMapCmd)

	for _, playerID := range playerIDs {
		key := fmt.Sprintf("%s%d", PlayerDetailsPrefix, playerID)
		cmds[playerID] = pipe.HGetAll(s.ctx, key)
	}

	_, err := pipe.Exec(s.ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	result := make(map[uint]*model.PlayerRanking)
	for playerID, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil || len(data) == 0 {
			continue
		}

		combatPower, _ := strconv.ParseInt(data["combat_power"], 10, 64)
		totalPoints, _ := strconv.Atoi(data["total_points"])
		wins, _ := strconv.Atoi(data["wins"])
		losses, _ := strconv.Atoi(data["losses"])
		winRate, _ := strconv.ParseFloat(data["win_rate"], 64)
		updatedAt, _ := time.Parse(time.RFC3339, data["updated_at"])

		result[playerID] = &model.PlayerRanking{
			PlayerID:    playerID,
			Username:    data["username"],
			CombatPower: combatPower,
			TotalPoints: totalPoints,
			Wins:        wins,
			Losses:      losses,
			WinRate:     winRate,
			UpdatedAt:   updatedAt,
		}
	}

	return result, nil
}

// GetThreshold returns the current top 10K threshold
func (s *LeaderboardService) GetThreshold() (int64, error) {
	val, err := s.redis.Get(s.ctx, ThresholdKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// updateThreshold updates the threshold to the 10,000th player's score
func (s *LeaderboardService) updateThreshold() error {
	// Get the score of the 10,000th player (0-indexed: 9999)
	result, err := s.redis.ZRevRangeWithScores(s.ctx, LeaderboardKey, Top10KLimit-1, Top10KLimit-1).Result()
	if err != nil || len(result) == 0 {
		return err
	}

	threshold := int64(result[0].Score)
	return s.redis.Set(s.ctx, ThresholdKey, threshold, 0).Err()
}

// SetMetadata sets leaderboard metadata
func (s *LeaderboardService) SetMetadata(totalPlayers int64) error {
	pipe := s.redis.Pipeline()
	pipe.Set(s.ctx, TotalPlayersKey, totalPlayers, 0)
	pipe.Set(s.ctx, LastSyncKey, time.Now().Format(time.RFC3339), 0)
	_, err := pipe.Exec(s.ctx)
	return err
}

// GetMetadata retrieves leaderboard metadata
func (s *LeaderboardService) GetMetadata() (*model.LeaderboardMetadata, error) {
	pipe := s.redis.Pipeline()
	totalCmd := pipe.Get(s.ctx, TotalPlayersKey)
	thresholdCmd := pipe.Get(s.ctx, ThresholdKey)
	lastSyncCmd := pipe.Get(s.ctx, LastSyncKey)
	_, err := pipe.Exec(s.ctx)

	metadata := &model.LeaderboardMetadata{}

	if total, err := totalCmd.Int64(); err == nil {
		metadata.TotalPlayers = total
	}

	if threshold, err := thresholdCmd.Int64(); err == nil {
		metadata.Top10KThreshold = threshold
	}

	if lastSync, err := lastSyncCmd.Result(); err == nil {
		metadata.LastUpdated, _ = time.Parse(time.RFC3339, lastSync)
	}

	return metadata, err
}

// AcquireSyncLock acquires a distributed lock for sync operations
func (s *LeaderboardService) AcquireSyncLock() (bool, error) {
	return s.redis.SetNX(s.ctx, SyncLockKey, s.instanceID, LockTimeout).Result()
}

// ReleaseSyncLock releases the distributed lock
func (s *LeaderboardService) ReleaseSyncLock() error {
	// Only release if we own the lock
	val, err := s.redis.Get(s.ctx, SyncLockKey).Result()
	if err != nil || val != s.instanceID {
		return err
	}
	return s.redis.Del(s.ctx, SyncLockKey).Err()
}

// ClearLeaderboard clears the entire leaderboard (use with caution)
func (s *LeaderboardService) ClearLeaderboard() error {
	return s.redis.Del(s.ctx, LeaderboardKey).Err()
}

// GetLeaderboardSize returns the current size of the leaderboard
func (s *LeaderboardService) GetLeaderboardSize() (int64, error) {
	return s.redis.ZCard(s.ctx, LeaderboardKey).Result()
}

// PublishPowerUpdateEvent publishes a power update event (for pub/sub pattern)
func (s *LeaderboardService) PublishPowerUpdateEvent(event *model.PowerUpdateEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.redis.Publish(s.ctx, "power.updates", data).Err()
}
