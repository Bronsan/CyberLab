package handlers

import (
	"strconv"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	challengeService *services.ChallengeService
}

func NewChallengeHandler(challengeService *services.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{challengeService: challengeService}
}

// GetChallenges returns a paginated list of challenges
// @Summary List challenges
// @Tags Challenges
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param category query string false "Category filter"
// @Param difficulty query string false "Difficulty filter"
// @Param search query string false "Search query"
// @Success 200 {object} utils.Response{data=services.ChallengeListResponse}
// @Router /api/v1/challenges [get]
func (h *ChallengeHandler) GetChallenges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	category := c.Query("category")
	difficulty := c.Query("difficulty")
	search := c.Query("search")

	result, err := h.challengeService.GetList(page, pageSize, category, difficulty, search)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// GetChallenge returns a single challenge by ID
// @Summary Get challenge detail
// @Tags Challenges
// @Produce json
// @Param id path int true "Challenge ID"
// @Success 200 {object} utils.Response{data=models.Challenge}
// @Router /api/v1/challenges/{id} [get]
func (h *ChallengeHandler) GetChallenge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid challenge id")
		return
	}

	challenge, err := h.challengeService.GetDetail(uint(id))
	if err != nil {
		utils.Error(c, 40004, err.Error())
		return
	}

	utils.Success(c, challenge)
}

// SearchChallenges searches challenges by keyword
// @Summary Search challenges
// @Tags Challenges
// @Produce json
// @Param q query string true "Search keyword"
// @Success 200 {object} utils.Response{data=[]models.Challenge}
// @Router /api/v1/challenges/search [get]
func (h *ChallengeHandler) SearchChallenges(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.BadRequest(c, "search query is required")
		return
	}

	results, err := h.challengeService.SearchChallenges(query)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, results)
}

// GetCategories returns all challenge categories
// @Summary Get categories
// @Tags Challenges
// @Produce json
// @Success 200 {object} utils.Response{data=[]string}
// @Router /api/v1/challenges/categories [get]
func (h *ChallengeHandler) GetCategories(c *gin.Context) {
	categories, err := h.challengeService.GetCategories()
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}
	utils.Success(c, categories)
}

// CreateChallenge creates a new challenge (admin only)
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	var challenge models.Challenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		utils.BadRequest(c, "invalid input: "+err.Error())
		return
	}

	if err := h.challengeService.CreateChallenge(&challenge); err != nil {
		utils.ServerError(c, "failed to create challenge")
		return
	}

	utils.SuccessWithMessage(c, "challenge created", challenge)
}

// UpdateChallenge updates an existing challenge (admin only)
func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid challenge id")
		return
	}

	var challenge models.Challenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		utils.BadRequest(c, "invalid input: "+err.Error())
		return
	}
	challenge.ID = uint(id)

	if err := h.challengeService.UpdateChallenge(&challenge); err != nil {
		utils.ServerError(c, "failed to update challenge")
		return
	}

	utils.SuccessWithMessage(c, "challenge updated", nil)
}

// DeleteChallenge deletes a challenge (admin only)
func (h *ChallengeHandler) DeleteChallenge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid challenge id")
		return
	}

	if err := h.challengeService.DeleteChallenge(uint(id)); err != nil {
		utils.ServerError(c, "failed to delete challenge")
		return
	}

	utils.SuccessWithMessage(c, "challenge deleted", nil)
}
