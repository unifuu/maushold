package repository

import (
	"maushold/player-service/model"

	"gorm.io/gorm"
)

type PlayerMonsterRepository interface {
	Create(monster *model.PlayerMonster) error
	FindByPlayerID(playerID uint) ([]model.PlayerMonster, error)
	FindByID(id uint) (*model.PlayerMonster, error)
}

type playerMonsterRepository struct {
	db *gorm.DB
}

func NewPlayerMonsterRepository(db *gorm.DB) PlayerMonsterRepository {
	return &playerMonsterRepository{db: db}
}

func (r *playerMonsterRepository) Create(monster *model.PlayerMonster) error {
	return r.db.Create(monster).Error
}

func (r *playerMonsterRepository) FindByPlayerID(playerID uint) ([]model.PlayerMonster, error) {
	var monster []model.PlayerMonster
	err := r.db.Where("player_id = ?", playerID).Find(&monster).Error
	return monster, err
}

func (r *playerMonsterRepository) FindByID(id uint) (*model.PlayerMonster, error) {
	var monster model.PlayerMonster
	err := r.db.First(&monster, id).Error
	return &monster, err
}
