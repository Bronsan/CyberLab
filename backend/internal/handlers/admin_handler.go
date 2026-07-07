package handlers

import (
	"strconv"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	userRepo          *repository.UserRepository
	containerService  *services.ContainerService
	challengeService  *services.ChallengeService
	announcementRepo  *repository.AnnouncementRepository
	logRepo           *repository.LogRepository
}

func NewAdminHandler(
	userRepo *repository.UserRepository,
	containerService *services.ContainerService,
	challengeService *services.ChallengeService,
	announcementRepo *repository.AnnouncementRepository,
	logRepo *repository.LogRepository,
) *AdminHandler {
	return &AdminHandler{
		userRepo:         userRepo,
		containerService: containerService,
		challengeService: challengeService,
		announcementRepo: announcementRepo,
		logRepo:          logRepo,
	}
}

// Dashboard returns admin dashboard stats
func (h *AdminHandler) Dashboard(c *gin.Context) {
	userCount, _ := h.userRepo.Count()
	challengeCount, _ := h.challengeService.GetList(1, 1, "", "", "")
	containerStats, _ := h.containerService.GetContainerStats()

	// Count running containers
	runningCount := int64(0)
	if v, ok := containerStats["RUNNING"]; ok {
		runningCount = v
	}

	utils.Success(c, gin.H{
		"totalUsers":      userCount,
		"totalChallenges": challengeCount.Total,
		"runningContainers": runningCount,
	})
}

// GetContainers returns all container instances
func (h *AdminHandler) GetContainers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	instances, total, err := h.containerService.GetAllInstances(page, pageSize)
	if err != nil {
		utils.ServerError(c, "failed to get instances")
		return
	}

	utils.Success(c, gin.H{
		"list":  instances,
		"total": total,
	})
}

// ForceDestroyContainer forcefully destroys a container
func (h *AdminHandler) ForceDestroyContainer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid instance id")
		return
	}

	if err := h.containerService.DestroyContainer(uint(id)); err != nil {
		utils.ServerError(c, "failed to destroy container")
		return
	}

	utils.SuccessWithMessage(c, "container destroyed", nil)
}

// CreateAnnouncement creates a new announcement
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		utils.BadRequest(c, "invalid input")
		return
	}

	if err := h.announcementRepo.Create(&announcement); err != nil {
		utils.ServerError(c, "failed to create announcement")
		return
	}

	utils.SuccessWithMessage(c, "announcement created", nil)
}

// DeleteAnnouncement deletes an announcement
func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid announcement id")
		return
	}

	if err := h.announcementRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "failed to delete announcement")
		return
	}

	utils.SuccessWithMessage(c, "announcement deleted", nil)
}

// GetUsers returns all users (admin)
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	users, total, err := h.userRepo.FindAll(page, pageSize)
	if err != nil {
		utils.ServerError(c, "failed to get users")
		return
	}

	utils.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

// GetLogs returns system logs
func (h *AdminHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	logs, total, err := h.logRepo.FindAll(page, pageSize)
	if err != nil {
		utils.ServerError(c, "failed to get logs")
		return
	}

	utils.Success(c, gin.H{
		"list":  logs,
		"total": total,
	})
}
