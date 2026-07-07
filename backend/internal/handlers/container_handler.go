package handlers

import (
	"strconv"
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	containerService *services.ContainerService
}

func NewContainerHandler(containerService *services.ContainerService) *ContainerHandler {
	return &ContainerHandler{containerService: containerService}
}

// StartChallenge starts a new container for a challenge
// @Summary Start challenge container
// @Tags Container
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object{challengeId=int} true "Challenge ID"
// @Success 200 {object} utils.Response{data=services.StartContainerResponse}
// @Router /api/v1/container/start [post]
func (h *ContainerHandler) StartChallenge(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req struct {
		ChallengeID uint `json:"challengeId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "challengeId is required")
		return
	}

	result, err := h.containerService.StartContainer(userID.(uint), req.ChallengeID)
	if err != nil {
		utils.Error(c, 50003, err.Error())
		return
	}

	utils.Success(c, result)
}

// GetContainerStatus returns the status of a container
// @Summary Get container status
// @Tags Container
// @Security Bearer
// @Produce json
// @Param id path int true "Instance ID"
// @Success 200 {object} utils.Response{data=object}
// @Router /api/v1/container/status/{id} [get]
func (h *ContainerHandler) GetContainerStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid instance id")
		return
	}

	instance, err := h.containerService.GetContainerStatus(uint(id))
	if err != nil {
		utils.Error(c, 40004, err.Error())
		return
	}

	runningTime := ""
	if instance.Status == models.InstanceStatusRunning && !instance.StartTime.IsZero() {
		runningTime = time.Since(instance.StartTime).Round(time.Minute).String()
	}

	utils.Success(c, gin.H{
		"status":      instance.Status,
		"runningTime": runningTime,
		"hostPort":    instance.HostPort,
	})
}

// StopChallenge stops a running container
// @Summary Stop container
// @Tags Container
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body object{instanceId=int} true "Instance ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/container/stop [post]
func (h *ContainerHandler) StopChallenge(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req struct {
		InstanceID uint `json:"instanceId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "instanceId is required")
		return
	}

	if err := h.containerService.StopContainer(userID.(uint), req.InstanceID); err != nil {
		utils.Error(c, 50003, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "container stopped", nil)
}

// GetUserInstances returns all running containers for the current user
func (h *ContainerHandler) GetUserInstances(c *gin.Context) {
	userID, _ := c.Get("userId")

	instances, err := h.containerService.GetUserInstances(userID.(uint))
	if err != nil {
		utils.ServerError(c, "failed to get instances")
		return
	}

	utils.Success(c, instances)
}
