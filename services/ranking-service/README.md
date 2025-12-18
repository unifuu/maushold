# Ranking Service

## PostgreSQL

``` sql
-- Index for power sorting
CREATE INDEX idx_players_power ON players(power DESC, id);

-- Index for user_id lookup
CREATE INDEX idx_players_user_id ON players(user_id);

-- Partitioning by power range
CREATE TABLE players_top (
    CHECK (power >= 1000000)
) INHERITS (players)

CREATE TABLE players_mid (
    CHECK (power >= 100000 AND power < 1000000)
) INHERITS (players);

CREATE TABLE players_low (
    CHECK (power < 100000)
) INHERITS (players);
```

## Redis

``` redis
# Main leaderboard - stores top 10,000 players
ZADD leaderboard:global {score} {user_id}

# Score = power value
# Member = user_id

Example:
ZADD leaderboard:global 999999 "player_123"
ZADD leaderboard:global 888888 "player_456"

# Get top 100
ZREVRANGE leaderboard:global 0 99 WITHSCORES

# Get player rank
ZREVRANK leaderboard:global "player_123"

# Get player's neighbors (nearby players)
ZREVRANGE leaderboard:global {rank-5} {rank+5} WITHSCORES

# Cache player details
HSET player:123 username "PlayerName" power "999999" level "50" avatar "url"

# Bulk get multiple players
HMGET player:123 username power level

# Last update timestamp
SET leaderboard:last_update "2024-12-10T10:30:00Z"

# Total players count
SET leaderboard:total_players "10000000"

# Update in progress flag
SET leaderboard:updating "true" EX 300
```