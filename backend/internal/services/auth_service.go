package services

import (
	"errors"
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	logRepo     *repository.LogRepository
	jwtSecret   string
	expireHours int
}

func NewAuthService(userRepo *repository.UserRepository, logRepo *repository.LogRepository, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		logRepo:     logRepo,
		jwtSecret:   jwtSecret,
		expireHours: expireHours,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type UserProfile struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
	Role        string `json:"role"`
	Score       int    `json:"score"`
	SolvedCount int    `json:"solvedCount"`
}

func (s *AuthService) Register(req *RegisterRequest, ip string) (*UserProfile, error) {
	// Check if user already exists
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already registered")
	}
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		Score:        0,
		SolvedCount:  0,
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Log
	s.logRepo.Create(&models.SystemLog{
		UserID: user.ID,
		Action: models.ActionRegister,
		IP:     ip,
	})

	return &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Score:       user.Score,
		SolvedCount: user.SolvedCount,
	}, nil
}

func (s *AuthService) Login(req *LoginRequest, ip string) (*LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status == 0 {
		return nil, errors.New("account has been banned")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret, s.expireHours)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login
	s.userRepo.UpdateFields(user.ID, map[string]interface{}{
		"updated_at": time.Now(),
	})

	// Log
	s.logRepo.Create(&models.SystemLog{
		UserID: user.ID,
		Action: models.ActionLogin,
		IP:     ip,
	})

	return &LoginResponse{
		Token: token,
		User: UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			Role:        user.Role,
			Score:       user.Score,
			SolvedCount: user.SolvedCount,
		},
	}, nil
}

func (s *AuthService) GetProfile(userID uint) (*UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &UserProfile{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Role:        user.Role,
		Score:       user.Score,
		SolvedCount: user.SolvedCount,
	}, nil
}

func (s *AuthService) UpdateProfile(userID uint, avatar, bio string) error {
	updates := make(map[string]interface{})
	if avatar != "" {
		updates["avatar"] = avatar
	}
	if bio != "" {
		updates["bio"] = bio
	}
	updates["updated_at"] = time.Now()

	return s.userRepo.UpdateFields(userID, updates)
}
