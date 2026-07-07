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
	rankingService *services.RankingService
	logRepo        *repository.LogRepository
}

func NewSubmissionHandler(
	submissionRepo *repository.SubmissionRepository,
	challengeRepo *repository.ChallengeRepository,
	rankingService *services.RankingService,
	logRepo *repository.LogRepository,
) *SubmissionHandler {
	return &SubmissionHandler{
		submissionRepo: submissionRepo,
		challengeRepo:  challengeRepo,
		rankingService: rankingService,
		logRepo:        logRepo,
	}
}

// SubmitFlag handles flag submission
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

	// Get challenge
	challenge, err := h.challengeRepo.FindByID(req.ChallengeID)
	if err != nil {
		utils.Error(c, 40004, "challenge not found")
		return
	}

	// Check if already solved
	_, err = h.submissionRepo.FindByUserAndChallenge(userID.(uint), req.ChallengeID)
	if err == nil {
		utils.Error(c, 200, "already solved")
		return
	}

	// Verify flag
	isCorrect := req.Flag == challenge.Flag

	// Record submission
	submission := &models.Submission{
		UserID:        userID.(uint),
		ChallengeID:   req.ChallengeID,
		SubmittedFlag: req.Flag,
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
// @Summary Get submission history
// @Tags Submission
// @Security Bearer
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} utils.Response
// @Router /api/v1/submit/history [get]
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
