package handlers

import (
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	rankingService *services.RankingService
}

func NewRankingHandler(rankingService *services.RankingService) *RankingHandler {
	return &RankingHandler{rankingService: rankingService}
}

// GetGlobalRanking returns the global ranking
// @Summary Get global ranking
// @Tags Ranking
// @Produce json
// @Success 200 {object} utils.Response{data=services.RankingResponse}
// @Router /api/v1/ranking/global [get]
func (h *RankingHandler) GetGlobalRanking(c *gin.Context) {
	ranking, err := h.rankingService.GetGlobalRanking()
	if err != nil {
		utils.ServerError(c, "failed to get ranking")
		return
	}
	utils.Success(c, ranking)
}

// GetWeeklyRanking returns the weekly ranking
// @Summary Get weekly ranking
// @Tags Ranking
// @Produce json
// @Success 200 {object} utils.Response{data=services.RankingResponse}
// @Router /api/v1/ranking/week [get]
func (h *RankingHandler) GetWeeklyRanking(c *gin.Context) {
	ranking, err := h.rankingService.GetWeeklyRanking()
	if err != nil {
		utils.ServerError(c, "failed to get weekly ranking")
		return
	}
	utils.Success(c, ranking)
}

// GetMonthlyRanking returns the monthly ranking
// @Summary Get monthly ranking
// @Tags Ranking
// @Produce json
// @Success 200 {object} utils.Response{data=services.RankingResponse}
// @Router /api/v1/ranking/month [get]
func (h *RankingHandler) GetMonthlyRanking(c *gin.Context) {
	ranking, err := h.rankingService.GetMonthlyRanking()
	if err != nil {
		utils.ServerError(c, "failed to get monthly ranking")
		return
	}
	utils.Success(c, ranking)
}
