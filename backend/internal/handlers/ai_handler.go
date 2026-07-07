package handlers

import (
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService *services.AIService
}

func NewAIHandler(aiService *services.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// GetHint returns an AI-generated hint for a challenge
// @Summary Get AI hint
// @Tags AI
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body services.HintRequest true "Hint request"
// @Success 200 {object} utils.Response{data=services.HintResponse}
// @Router /api/v1/ai/hint [post]
func (h *AIHandler) GetHint(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req services.HintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid input: "+err.Error())
		return
	}

	result, err := h.aiService.GetHint(userID.(uint), &req)
	if err != nil {
		utils.Error(c, 50004, err.Error())
		return
	}

	utils.Success(c, result)
}

// AuditCode performs AI code audit
// @Summary AI code audit
// @Tags AI
// @Security Bearer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Source code file (zip/rar/tar.gz)"
// @Success 200 {object} utils.Response{data=services.AuditResponse}
// @Router /api/v1/ai/audit [post]
func (h *AIHandler) AuditCode(c *gin.Context) {
	userID, _ := c.Get("userId")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	buf := make([]byte, 1024*1024) // 1MB max
	n, err := file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		utils.ServerError(c, "failed to read file")
		return
	}

	content := string(buf[:n])
	result, err := h.aiService.AuditCode(userID.(uint), content, header.Filename)
	if err != nil {
		utils.Error(c, 50004, err.Error())
		return
	}

	utils.Success(c, result)
}
