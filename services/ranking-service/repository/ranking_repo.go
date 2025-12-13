package repository

import (
	"maushold/ranking-service/model"

	"gorm.io/gorm"
)

type RankingRepository interface {
	Create(ranking *model.PlayerRanking) error
	Update(ranking *model.PlayerRanking) error
	FindByPlayerID(playerID uint) (*model.PlayerRanking, error)
	FindAll() ([]model.PlayerRanking, error)
	FindTopN(limit int) ([]model.PlayerRanking, error)
}

type rankingRepository struct {
	db *gorm.DB
}

func NewRankingRepository(db *gorm.DB) RankingRepository {
	return &rankingRepository{db: db}
}

func (r *rankingRepository) Create(ranking *model.PlayerRanking) error {
	return r.db.Create(ranking).Error
}

func (r *rankingRepository) Update(ranking *model.PlayerRanking) error {
	return r.db.Save(ranking).Error
}

func (r *rankingRepository) FindByPlayerID(playerID uint) (*model.PlayerRanking, error) {
	var ranking model.PlayerRanking
	err := r.db.Where("player_id = ?", playerID).First(&ranking).Error
	return &ranking, err
}

func (r *rankingRepository) FindAll() ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Order("total_points DESC").Find(&rankings).Error
	return rankings, err
}

func (r *rankingRepository) FindTopN(limit int) ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Order("total_points DESC").Limit(limit).Find(&rankings).Error
	return rankings, err
}
