package services

import (
	"errors"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
)

type ChallengeService struct {
	challengeRepo *repository.ChallengeRepository
}

func NewChallengeService(challengeRepo *repository.ChallengeRepository) *ChallengeService {
	return &ChallengeService{challengeRepo: challengeRepo}
}

type ChallengeListResponse struct {
	List  []models.Challenge `json:"list"`
	Total int64              `json:"total"`
}

func (s *ChallengeService) GetList(page, pageSize int, category, difficulty, search string) (*ChallengeListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	challenges, total, err := s.challengeRepo.FindAll(page, pageSize, category, difficulty, search)
	if err != nil {
		return nil, errors.New("failed to fetch challenges")
	}

	return &ChallengeListResponse{
		List:  challenges,
		Total: total,
	}, nil
}

func (s *ChallengeService) GetDetail(id uint) (*models.Challenge, error) {
	challenge, err := s.challengeRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("challenge not found")
	}
	return challenge, nil
}

func (s *ChallengeService) SearchChallenges(query string) ([]models.Challenge, error) {
	challenges, _, err := s.challengeRepo.FindAll(1, 20, "", "", query)
	if err != nil {
		return nil, errors.New("search failed")
	}
	return challenges, nil
}

func (s *ChallengeService) CreateChallenge(challenge *models.Challenge) error {
	return s.challengeRepo.Create(challenge)
}

func (s *ChallengeService) UpdateChallenge(challenge *models.Challenge) error {
	return s.challengeRepo.Update(challenge)
}

func (s *ChallengeService) DeleteChallenge(id uint) error {
	return s.challengeRepo.Delete(id)
}

func (s *ChallengeService) GetCategories() ([]string, error) {
	return s.challengeRepo.GetCategories()
}
