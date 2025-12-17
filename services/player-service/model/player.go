package model

import "time"

type Player struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // "-" means don't include in JSON responses
	Points    int       `gorm:"default:0" json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlayerMonster struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlayerID   uint      `gorm:"not null;index" json:"player_id"`
	MonsterID  int       `gorm:"not null" json:"monster_id"`
	Nickname   string    `gorm:"size:255" json:"nickname"`
	Level      int       `gorm:"default:1" json:"level"`
	Experience int       `gorm:"default:0" json:"experience"`
	HP         int       `gorm:"not null" json:"hp"`
	Attack     int       `gorm:"not null" json:"attack"`
	Defense    int       `gorm:"not null" json:"defense"`
	Speed      int       `gorm:"not null" json:"speed"`
	CreatedAt  time.Time `json:"created_at"`
}
