package handlers

import (
	"strconv"

	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementRepo *repository.AnnouncementRepository
}

func NewAnnouncementHandler(announcementRepo *repository.AnnouncementRepository) *AnnouncementHandler {
	return &AnnouncementHandler{announcementRepo: announcementRepo}
}

// GetAnnouncements returns a list of announcements
// @Summary Get announcements
// @Tags Announcements
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} utils.Response
// @Router /api/v1/announcement [get]
func (h *AnnouncementHandler) GetAnnouncements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	announcements, total, err := h.announcementRepo.FindAll(page, pageSize)
	if err != nil {
		utils.ServerError(c, "failed to get announcements")
		return
	}

	utils.Success(c, gin.H{
		"list":  announcements,
		"total": total,
	})
}

// GetAnnouncement returns a single announcement by ID
// @Summary Get announcement detail
// @Tags Announcements
// @Produce json
// @Param id path int true "Announcement ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/announcement/{id} [get]
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid announcement id")
		return
	}

	announcement, err := h.announcementRepo.FindByID(uint(id))
	if err != nil {
		utils.Error(c, 40004, "announcement not found")
		return
	}

	utils.Success(c, announcement)
}
