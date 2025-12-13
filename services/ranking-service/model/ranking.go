package model

import "time"

type PlayerRanking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PlayerID     uint      `gorm:"unique;not null;index" json:"player_id"`
	Username     string    `json:"username"`
	TotalPoints  int       `gorm:"default:1000;index" json:"total_points"`
	TotalBattles int       `gorm:"default:0" json:"total_battles"`
	Wins         int       `gorm:"default:0" json:"wins"`
	Losses       int       `gorm:"default:0" json:"losses"`
	WinRate      float64   `json:"win_rate"`
	Rank         int       `json:"rank"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LeaderboardEntry struct {
	PlayerID    uint    `json:"player_id"`
	Username    string  `json:"username"`
	TotalPoints int     `json:"total_points"`
	Wins        int     `json:"wins"`
	Losses      int     `json:"losses"`
	WinRate     float64 `json:"win_rate"`
	Rank        int     `json:"rank"`
}

type Player struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Points   int    `json:"points"`
}
