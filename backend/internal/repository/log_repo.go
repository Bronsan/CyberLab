package repository

import (
	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) Create(log *models.SystemLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepository) FindByUser(userID uint, page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	query := r.db.Model(&models.SystemLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepository) FindAll(page, pageSize int) ([]models.SystemLog, int64, error) {
	var logs []models.SystemLog
	var total int64

	r.db.Model(&models.SystemLog{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&logs).Error
	return logs, total, err
}
