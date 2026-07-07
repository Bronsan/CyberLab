package repository

import (
	"time"

	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type InstanceRepository struct {
	db *gorm.DB
}

func NewInstanceRepository(db *gorm.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

func (r *InstanceRepository) Create(instance *models.ChallengeInstance) error {
	return r.db.Create(instance).Error
}

func (r *InstanceRepository) FindByID(id uint) (*models.ChallengeInstance, error) {
	var instance models.ChallengeInstance
	err := r.db.Preload("User").Preload("Challenge").First(&instance, id).Error
	return &instance, err
}

func (r *InstanceRepository) FindByUserAndChallenge(userID, challengeID uint) (*models.ChallengeInstance, error) {
	var instance models.ChallengeInstance
	err := r.db.Where("user_id = ? AND challenge_id = ? AND status IN ?",
		userID, challengeID, []string{models.InstanceStatusCreated, models.InstanceStatusRunning}).
		First(&instance).Error
	return &instance, err
}

func (r *InstanceRepository) FindRunningByUser(userID uint) ([]models.ChallengeInstance, error) {
	var instances []models.ChallengeInstance
	err := r.db.Where("user_id = ? AND status = ?", userID, models.InstanceStatusRunning).
		Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.ChallengeInstance{}).Where("id = ?", id).Update("status", status).Error
}

func (r *InstanceRepository) Update(instance *models.ChallengeInstance) error {
	return r.db.Save(instance).Error
}

func (r *InstanceRepository) FindExpired() ([]models.ChallengeInstance, error) {
	var instances []models.ChallengeInstance
	err := r.db.Where("status IN ? AND expire_time < ?",
		[]string{models.InstanceStatusRunning, models.InstanceStatusCreated},
		time.Now()).
		Find(&instances).Error
	return instances, err
}

func (r *InstanceRepository) FindAll(page, pageSize int) ([]models.ChallengeInstance, int64, error) {
	var instances []models.ChallengeInstance
	var total int64

	r.db.Model(&models.ChallengeInstance{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Preload("User").Preload("Challenge").
		Offset(offset).Limit(pageSize).Order("id desc").Find(&instances).Error
	return instances, total, err
}

func (r *InstanceRepository) CountRunning() (int64, error) {
	var count int64
	err := r.db.Model(&models.ChallengeInstance{}).Where("status = ?", models.InstanceStatusRunning).Count(&count).Error
	return count, err
}

func (r *InstanceRepository) CountByStatus() (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}
	err := r.db.Model(&models.ChallengeInstance{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, r := range results {
		counts[r.Status] = r.Count
	}
	return counts, nil
}
