package model

import "time"

type Player struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Points    int       `gorm:"default:0" json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlayerPokemon struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PlayerID   uint      `gorm:"not null;index" json:"player_id"`
	PokemonID  int       `gorm:"not null" json:"pokemon_id"`
	Level      int       `gorm:"default:1" json:"level"`
	Experience int       `gorm:"default:0" json:"experience"`
	HP         int       `gorm:"not null" json:"hp"`
	Attack     int       `gorm:"not null" json:"attack"`
	Defense    int       `gorm:"not null" json:"defense"`
	Speed      int       `gorm:"not null" json:"speed"`
	CreatedAt  time.Time `json:"created_at"`
}
