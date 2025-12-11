package repository

import (
	"player-service/model"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(player *model.Player) error
	FindByID(id uint) (*model.Player, error)
	Update(player *model.Player) error
	FindAll() ([]model.Player, error)
	UpdatePoints(id uint, points int) error
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *model.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) FindByID(id uint) (*model.Player, error) {
	var player model.Player
	err := r.db.First(&player, id).Error
	return &player, err
}

func (r *playerRepository) Update(player *model.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) FindAll() ([]model.Player, error) {
	var players []model.Player
	err := r.db.Order("points DESC").Find(&players).Error
	return players, err
}

func (r *playerRepository) UpdatePoints(id uint, points int) error {
	return r.db.Model(&model.Player{}).Where("id = ?", id).Update("points", points).Error
}
