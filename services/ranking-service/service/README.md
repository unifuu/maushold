## `leaderboard_srv.go`

This is a low-level service that specializes in Redis operations. It doesn't care about business rules (like how much power a player gained from a battle); it only cares about efficiently storing and retrieving data in-memory.

- Sorted Set Management: Handles the actual ZADD, ZREM, and ZREVRANK commands that manage the top 10,000 players.
- Hash Caching: Handles caching full player details (username, wins, losses) so we don't hit the database every time we view the leaderboard.
- Infrastructure Logic: Implements distributed locking (to prevent two servers from syncing at once) and the "Top 10K Threshold" logic to keep Redis slim.


### `GetThreshold()`

``` go
// GetThreshold returns the current top 10K threshold
func (s *LeaderboardService) GetThreshold() (int64, error) {
	val, err := s.redis.Get(s.ctx, ThresholdKey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}
```

What is the "Threshold"?
Simply put, the threshold is the Combat Power of the #10,000 player.

In a system with tens of millions of players, you don't want to store everyone in Redis. It would be too expensive and slow. Instead, we only keep the Top 10,000 in Redis.

How it works:
1. Filtering Updates: When a player's power changes, we call GetThreshold(). If their new power is lower than the threshold (meaning they aren't even close to the top 10,000), we don't update Redis at all. We only update the database. This saves 99% of unnecessary Redis writes.
2. The "Elite" Cache: If their power is higher than the threshold, we add them to Redis and then "trim" the leaderboard back down to 10,000.
3. Updating the Bar: After trimming, the current #10,000 player's score becomes the new threshold. This is what updateThreshold() does.

## `ranking_srv.go`

This is the high-level service that coordinates everyone else. It's the "brain" of the ranking system.

- Battle Logic: When a player wins, this service calculates their new combat power and decides which database and cache updates are needed.
- Tiered Caching: This is where I implemented the Redis-first logic you asked for. It says: "Go ask the 
LeaderboardService
 for the rank. If it fails, go ask the Database Repository."
- Sync Orchestration: It schedules when to refresh the database's Materialized View and when to push a fresh copy of everything into Redis.