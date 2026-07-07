package handlers

import (
	"strconv"
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SubmissionHandler struct {
	submissionRepo *repository.SubmissionRepository
	challengeRepo  *repository.ChallengeRepository
	instanceRepo   *repository.InstanceRepository
	rankingService *services.RankingService
	logRepo        *repository.LogRepository
	flagHashSalt   string
}

func NewSubmissionHandler(
	submissionRepo *repository.SubmissionRepository,
	challengeRepo *repository.ChallengeRepository,
	rankingService *services.RankingService,
	logRepo *repository.LogRepository,
	instanceRepo *repository.InstanceRepository,
	flagHashSalt string,
) *SubmissionHandler {
	if flagHashSalt == "" {
		flagHashSalt = "cyberlab-default-salt-change-in-production"
		utils.Log.Warn("FLAG_HASH_SALT not set, using default — set FLAG_HASH_SALT env var in production")
	}
	return &SubmissionHandler{
		submissionRepo: submissionRepo,
		challengeRepo:  challengeRepo,
		instanceRepo:   instanceRepo,
		rankingService: rankingService,
		logRepo:        logRepo,
		flagHashSalt:   flagHashSalt,
	}
}

// SubmitFlag handles flag submission with security validation.
// @Summary Submit a flag
// @Tags Submission
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object{challengeId=int,flag=string} true "Flag submission"
// @Success 200 {object} utils.Response
// @Router /api/v1/submit [post]
func (h *SubmissionHandler) SubmitFlag(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req struct {
		ChallengeID uint   `json:"challengeId" binding:"required"`
		Flag        string `json:"flag" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "challengeId and flag are required")
		return
	}

	// Basic input sanitization
	if len(req.Flag) > 1024 {
		utils.BadRequest(c, "flag too long")
		return
	}

	// Get challenge
	challenge, err := h.challengeRepo.FindByID(req.ChallengeID)
	if err != nil {
		utils.Error(c, 40004, "challenge not found")
		return
	}

	if !challenge.IsActive {
		utils.Error(c, 40004, "challenge is not active")
		return
	}

	// Verify the user has an active or previously assigned instance for this challenge
	// This prevents blind flag brute-forcing across challenges the user hasn't started
	instance, _ := h.instanceRepo.FindByUserAndChallenge(userID.(uint), req.ChallengeID)
	if instance == nil {
		// Check if user has ever submitted for this challenge (relaxed check)
		existing, _ := h.submissionRepo.FindByUserAndChallenge(userID.(uint), req.ChallengeID)
		if existing == nil {
			utils.Error(c, 40006, "you must start the challenge before submitting a flag")
			return
		}
	}

	// Check if already solved
	_, err = h.submissionRepo.FindByUserAndChallenge(userID.(uint), req.ChallengeID)
	if err == nil {
		utils.Success(c, gin.H{
			"correct": true,
			"score":   challenge.Score,
			"message": "already solved",
		})
		return
	}

	// Verify flag using hash comparison (not plaintext)
	isCorrect := utils.VerifyFlag(req.Flag, challenge.Flag, req.ChallengeID, h.flagHashSalt)

	// Record submission — store HASH of submitted flag, not plaintext
	submittedFlagHash := utils.HashString(req.Flag)
	submission := &models.Submission{
		UserID:        userID.(uint),
		ChallengeID:   req.ChallengeID,
		SubmittedFlag: submittedFlagHash, // store hash only
		IsCorrect:     isCorrect,
		SubmittedAt:   time.Now(),
	}
	h.submissionRepo.Create(submission)

	if !isCorrect {
		utils.Success(c, gin.H{
			"correct": false,
			"score":   0,
		})
		return
	}

	// Record solved challenge
	h.submissionRepo.CreateUserChallenge(&models.UserChallenge{
		UserID:      userID.(uint),
		ChallengeID: req.ChallengeID,
		Score:       challenge.Score,
		SolvedAt:    time.Now(),
	})

	// Update user score
	h.rankingService.UpdateUserScore(userID.(uint), challenge.Score)

	// Log
	h.logRepo.Create(&models.SystemLog{
		UserID: userID.(uint),
		Action: models.ActionSubmitFlag,
	})

	utils.Success(c, gin.H{
		"correct": true,
		"score":   challenge.Score,
	})
}

// GetSubmissionHistory returns the user's submission history
func (h *SubmissionHandler) GetSubmissionHistory(c *gin.Context) {
	userID, _ := c.Get("userId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	submissions, total, err := h.submissionRepo.FindByUser(userID.(uint), page, pageSize)
	if err != nil {
		utils.ServerError(c, "failed to get history")
		return
	}

	utils.Success(c, gin.H{
		"list":  submissions,
		"total": total,
	})
}

// GetSolvedChallenges returns the user's solved challenges
func (h *SubmissionHandler) GetSolvedChallenges(c *gin.Context) {
	userID, _ := c.Get("userId")

	solved, err := h.submissionRepo.GetSolvedChallenges(userID.(uint))
	if err != nil {
		utils.ServerError(c, "failed to get solved challenges")
		return
	}

	utils.Success(c, solved)
}
