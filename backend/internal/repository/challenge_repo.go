package repository

import (
	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}

func (r *ChallengeRepository) Create(challenge *models.Challenge) error {
	return r.db.Create(challenge).Error
}

func (r *ChallengeRepository) FindByID(id uint) (*models.Challenge, error) {
	var challenge models.Challenge
	err := r.db.Preload("Tags").First(&challenge, id).Error
	return &challenge, err
}

func (r *ChallengeRepository) FindAll(page, pageSize int, category, difficulty, search string) ([]models.Challenge, int64, error) {
	var challenges []models.Challenge
	var total int64

	query := r.db.Model(&models.Challenge{}).Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Tags").Offset(offset).Limit(pageSize).Order("id desc").Find(&challenges).Error
	return challenges, total, err
}

func (r *ChallengeRepository) Update(challenge *models.Challenge) error {
	return r.db.Save(challenge).Error
}

func (r *ChallengeRepository) Delete(id uint) error {
	r.db.Where("challenge_id = ?", id).Delete(&models.ChallengeTag{})
	return r.db.Delete(&models.Challenge{}, id).Error
}

func (r *ChallengeRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Challenge{}).Select("DISTINCT category").Pluck("category", &categories).Error
	return categories, err
}

func (r *ChallengeRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Challenge{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}
