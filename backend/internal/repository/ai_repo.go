package repository

import (
	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type AIRepository struct {
	db *gorm.DB
}

func NewAIRepository(db *gorm.DB) *AIRepository {
	return &AIRepository{db: db}
}

func (r *AIRepository) GetChallengeDetail(challengeID uint) (*models.Challenge, error) {
	var challenge models.Challenge
	err := r.db.Select("id, title, category, difficulty, description").First(&challenge, challengeID).Error
	return &challenge, err
}

func (r *AIRepository) CreateHint(hint *models.AIHint) error {
	return r.db.Create(hint).Error
}

func (r *AIRepository) GetHintHistory(userID uint, page, pageSize int) ([]models.AIHint, int64, error) {
	var hints []models.AIHint
	var total int64

	query := r.db.Model(&models.AIHint{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Challenge").Offset(offset).Limit(pageSize).Order("created_at desc").Find(&hints).Error
	return hints, total, err
}
