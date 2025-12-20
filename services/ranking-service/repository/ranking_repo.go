package repository

import (
	"maushold/ranking-service/model"
	"time"

	"gorm.io/gorm"
)

type RankingRepository interface {
	Create(ranking *model.PlayerRanking) error
	Update(ranking *model.PlayerRanking) error
	UpdateCombatPower(playerID uint, combatPower int64) error
	FindByPlayerID(playerID uint) (*model.PlayerRanking, error)
	FindByPlayerIDs(playerIDs []uint) ([]model.PlayerRanking, error)
	FindAll() ([]model.PlayerRanking, error)
	FindTopN(limit int) ([]model.PlayerRanking, error)
	FindTopNByCombatPower(limit int) ([]model.PlayerRanking, error)
	GetTop10KThreshold() (int64, error)
	GetTotalPlayerCount() (int64, error)
	RefreshMaterializedView() error
	FindTopNFromMaterializedView(limit int) ([]model.LeaderboardEntry, error)
	GetPlayerRankFromMaterializedView(playerID uint) (int, error)
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

func (r *rankingRepository) UpdateCombatPower(playerID uint, combatPower int64) error {
	return r.db.Model(&model.PlayerRanking{}).
		Where("player_id = ?", playerID).
		Updates(map[string]interface{}{
			"combat_power":   combatPower,
			"last_battle_at": time.Now(),
			"updated_at":     time.Now(),
		}).Error
}

func (r *rankingRepository) FindByPlayerID(playerID uint) (*model.PlayerRanking, error) {
	var ranking model.PlayerRanking
	err := r.db.Where("player_id = ?", playerID).First(&ranking).Error
	return &ranking, err
}

func (r *rankingRepository) FindByPlayerIDs(playerIDs []uint) ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Where("player_id IN ?", playerIDs).Find(&rankings).Error
	return rankings, err
}

func (r *rankingRepository) FindAll() ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Order("combat_power DESC, id ASC").Find(&rankings).Error
	return rankings, err
}

func (r *rankingRepository) FindTopN(limit int) ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Order("total_points DESC").Limit(limit).Find(&rankings).Error
	return rankings, err
}

func (r *rankingRepository) FindTopNByCombatPower(limit int) ([]model.PlayerRanking, error) {
	var rankings []model.PlayerRanking
	err := r.db.Where("total_points > 0").Order("combat_power DESC, id ASC").Limit(limit).Find(&rankings).Error
	return rankings, err
}

func (r *rankingRepository) GetTop10KThreshold() (int64, error) {
	var threshold int64
	err := r.db.Model(&model.PlayerRanking{}).
		Select("combat_power").
		Order("combat_power DESC").
		Limit(1).
		Offset(9999).
		Scan(&threshold).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return threshold, err
}

func (r *rankingRepository) GetTotalPlayerCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.PlayerRanking{}).Count(&count).Error
	return count, err
}

func (r *rankingRepository) RefreshMaterializedView() error {
	return r.db.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY top_10k_players").Error
}

func (r *rankingRepository) FindTopNFromMaterializedView(limit int) ([]model.LeaderboardEntry, error) {
	var entries []model.LeaderboardEntry
	err := r.db.Table("top_10k_players").
		Select("player_id, username, combat_power, total_points, wins, losses, win_rate, rank, updated_at").
		Order("rank ASC").
		Limit(limit).
		Scan(&entries).Error
	return entries, err
}

func (r *rankingRepository) GetPlayerRankFromMaterializedView(playerID uint) (int, error) {
	var rank int
	err := r.db.Table("top_10k_players").
		Select("rank").
		Where("player_id = ?", playerID).
		Scan(&rank).Error

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return rank, err
}
