package model

import "time"

type Pokemon struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Type1       string    `gorm:"not null" json:"type1"`
	Type2       string    `json:"type2"`
	BaseHP      int       `gorm:"not null" json:"base_hp"`
	BaseAttack  int       `gorm:"not null" json:"base_attack"`
	BaseDefense int       `gorm:"not null" json:"base_defense"`
	BaseSpeed   int       `gorm:"not null" json:"base_speed"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
}
