package repository

import (
	"maushold/battle-service/model"

	"gorm.io/gorm"
)

type BattleRepository interface {
	Create(battle *model.Battle) error
	FindByID(id uint) (*model.Battle, error)
	FindByPlayerID(playerID uint) ([]model.Battle, error)
	FindRecent(limit int) ([]model.Battle, error)
	Update(battle *model.Battle) error
}

type battleRepository struct {
	db *gorm.DB
}

func NewBattleRepository(db *gorm.DB) BattleRepository {
	return &battleRepository{db: db}
}

func (r *battleRepository) Create(battle *model.Battle) error {
	return r.db.Create(battle).Error
}

func (r *battleRepository) FindByID(id uint) (*model.Battle, error) {
	var battle model.Battle
	err := r.db.First(&battle, id).Error
	return &battle, err
}

func (r *battleRepository) FindByPlayerID(playerID uint) ([]model.Battle, error) {
	var battles []model.Battle
	err := r.db.Where("player1_id = ? OR player2_id = ?", playerID, playerID).
		Order("created_at DESC").
		Limit(20).
		Find(&battles).Error
	return battles, err
}

func (r *battleRepository) FindRecent(limit int) ([]model.Battle, error) {
	var battles []model.Battle
	err := r.db.Order("created_at DESC").Limit(limit).Find(&battles).Error
	return battles, err
}

func (r *battleRepository) Update(battle *model.Battle) error {
	return r.db.Save(battle).Error
}
