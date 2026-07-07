package repository

import (
	"github.com/cyberlab/backend/internal/models"
	"gorm.io/gorm"
)

type SubmissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(submission *models.Submission) error {
	return r.db.Create(submission).Error
}

func (r *SubmissionRepository) FindByUserAndChallenge(userID, challengeID uint) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.Where("user_id = ? AND challenge_id = ? AND is_correct = ?",
		userID, challengeID, true).First(&submission).Error
	return &submission, err
}

func (r *SubmissionRepository) FindByUser(userID uint, page, pageSize int) ([]models.Submission, int64, error) {
	var submissions []models.Submission
	var total int64

	query := r.db.Model(&models.Submission{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Challenge").Offset(offset).Limit(pageSize).Order("submitted_at desc").Find(&submissions).Error
	return submissions, total, err
}

func (r *SubmissionRepository) CreateUserChallenge(uc *models.UserChallenge) error {
	return r.db.Create(uc).Error
}

func (r *SubmissionRepository) FindUserChallenge(userID, challengeID uint) (*models.UserChallenge, error) {
	var uc models.UserChallenge
	err := r.db.Where("user_id = ? AND challenge_id = ?", userID, challengeID).First(&uc).Error
	return &uc, err
}

func (r *SubmissionRepository) GetSolvedChallenges(userID uint) ([]models.UserChallenge, error) {
	var solved []models.UserChallenge
	err := r.db.Where("user_id = ?", userID).Preload("Challenge").Find(&solved).Error
	return solved, err
}

func (r *SubmissionRepository) CountSolvedByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.UserChallenge{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
