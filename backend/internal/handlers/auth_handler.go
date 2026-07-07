package handlers

import (
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body services.RegisterRequest true "Registration info"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid input: "+err.Error())
		return
	}

	ip := c.ClientIP()
	user, err := h.authService.Register(&req, ip)
	if err != nil {
		utils.Error(c, 40004, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "register success", user)
}

// Login handles user login
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body services.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response{data=services.LoginResponse}
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid input: "+err.Error())
		return
	}

	ip := c.ClientIP()
	resp, err := h.authService.Login(&req, ip)
	if err != nil {
		utils.Error(c, 40005, err.Error())
		return
	}

	utils.Success(c, resp)
}

// GetProfile returns the current user's profile
// @Summary Get current user profile
// @Tags Auth
// @Security Bearer
// @Produce json
// @Success 200 {object} utils.Response{data=services.UserProfile}
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		utils.Error(c, 40004, err.Error())
		return
	}

	utils.Success(c, profile)
}

// UpdateProfile updates the current user's profile
// @Summary Update user profile
// @Tags Auth
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body map[string]string true "Profile fields to update"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")

	var req struct {
		Avatar string `json:"avatar"`
		Bio    string `json:"bio"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "invalid input")
		return
	}

	if err := h.authService.UpdateProfile(userID.(uint), req.Avatar, req.Bio); err != nil {
		utils.ServerError(c, "failed to update profile")
		return
	}

	utils.SuccessWithMessage(c, "profile updated", nil)
}
