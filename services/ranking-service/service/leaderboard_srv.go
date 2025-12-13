package service

import (
	"context"
	"fmt"
	"strconv"

	"maushold/ranking-service/model"

	"github.com/go-redis/redis/v8"
)

type LeaderboardService struct {
	redis *redis.Client
	ctx   context.Context
}

func NewLeaderboardService(redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *LeaderboardService) UpdatePlayerScore(playerID uint, points int) error {
	return s.redis.ZAdd(s.ctx, "leaderboard", &redis.Z{
		Score:  float64(points),
		Member: fmt.Sprintf("%d", playerID),
	}).Err()
}

func (s *LeaderboardService) GetTopPlayers(limit int) ([]model.LeaderboardEntry, error) {
	result, err := s.redis.ZRevRangeWithScores(s.ctx, "leaderboard", 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	entries := make([]model.LeaderboardEntry, 0, len(result))
	for _, z := range result {
		playerID, _ := strconv.ParseUint(z.Member.(string), 10, 64)
		entries = append(entries, model.LeaderboardEntry{
			PlayerID:    uint(playerID),
			TotalPoints: int(z.Score),
		})
	}

	return entries, nil
}

func (s *LeaderboardService) GetPlayerRank(playerID uint) (int64, error) {
	rank, err := s.redis.ZRevRank(s.ctx, "leaderboard", fmt.Sprintf("%d", playerID)).Result()
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (s *LeaderboardService) ClearLeaderboard() error {
	return s.redis.Del(s.ctx, "leaderboard").Err()
}
