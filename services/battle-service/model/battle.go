package model

import "time"

type Battle struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Player1ID   uint       `gorm:"not null;index" json:"player1_id"`
	Player2ID   uint       `gorm:"not null;index" json:"player2_id"`
	Monster1ID  uint       `gorm:"not null" json:"monster1_id"`
	Monster2ID  uint       `gorm:"not null" json:"monster2_id"`
	WinnerID    uint       `json:"winner_id"`
	Status      string     `gorm:"default:'pending'" json:"status"`
	BattleLog   string     `gorm:"type:text" json:"battle_log"`
	PointsWon   int        `json:"points_won"`
	PointsLost  int        `json:"points_lost"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type PlayerMonster struct {
	ID       uint   `json:"id"`
	PlayerID uint   `json:"player_id"`
	Nickname string `json:"nickname"`
	HP       int    `json:"hp"`
	Attack   int    `json:"attack"`
	Defense  int    `json:"defense"`
	Speed    int    `json:"speed"`
	Level    int    `json:"level"`
}
