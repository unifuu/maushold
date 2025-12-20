package model

import "time"

type PlayerRanking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlayerID     uint      `gorm:"unique;not null;index:idx_player_id" json:"player_id"`
	Username     string    `json:"username"`
	CombatPower  int64     `gorm:"default:0;index:idx_combat_power,sort:desc" json:"combat_power"`
	TotalPoints  int       `gorm:"default:1000" json:"total_points"`
	TotalBattles int       `gorm:"default:0" json:"total_battles"`
	Wins         int       `gorm:"default:0" json:"wins"`
	Losses       int       `gorm:"default:0" json:"losses"`
	WinRate      float64   `json:"win_rate"`
	Rank         int       `gorm:"-" json:"rank"` // Computed field, not stored
	LastBattleAt time.Time `json:"last_battle_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LeaderboardEntry struct {
	PlayerID    uint      `json:"player_id"`
	Username    string    `json:"username"`
	CombatPower int64     `json:"combat_power"`
	TotalPoints int       `json:"total_points"`
	Wins        int       `json:"wins"`
	Losses      int       `json:"losses"`
	WinRate     float64   `json:"win_rate"`
	Rank        int       `json:"rank"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LeaderboardMetadata struct {
	TotalPlayers    int64     `json:"total_players"`
	Top10KThreshold int64     `json:"top_10k_threshold"`
	LastUpdated     time.Time `json:"last_updated"`
	CacheHit        bool      `json:"cache_hit"`
}

type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry  `json:"leaderboard"`
	Metadata    LeaderboardMetadata `json:"metadata"`
}

type PlayerRankContext struct {
	Player    LeaderboardEntry    `json:"player"`
	Neighbors []LeaderboardEntry  `json:"neighbors"`
	Metadata  LeaderboardMetadata `json:"metadata"`
}

type Player struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Points   int    `json:"points"`
}

// PowerUpdateEvent represents a combat power update event for message queue
type PowerUpdateEvent struct {
	PlayerID    uint      `json:"player_id"`
	Username    string    `json:"username"`
	CombatPower int64     `json:"combat_power"`
	PointsDelta int       `json:"points_delta"`
	IsWin       bool      `json:"is_win"`
	Timestamp   time.Time `json:"timestamp"`
}

type Battle struct {
	ID        uint   `json:"id"`
	Player1ID uint   `json:"player1_id"`
	Player2ID uint   `json:"player2_id"`
	WinnerID  uint   `json:"winner_id"`
	Status    string `json:"status"`
}
