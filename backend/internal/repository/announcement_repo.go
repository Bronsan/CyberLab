package repository

import (
	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

func (r *AnnouncementRepository) Create(announcement *models.Announcement) error {
	return r.db.Create(announcement).Error
}

func (r *AnnouncementRepository) FindByID(id uint) (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.First(&announcement, id).Error
	return &announcement, err
}

func (r *AnnouncementRepository) FindAll(page, pageSize int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	r.db.Model(&models.Announcement{}).Count(&total)

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&announcements).Error
	return announcements, total, err
}

func (r *AnnouncementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Announcement{}, id).Error
}
