package router

import (
	"fmt"

	"github.com/cyberlab/backend/internal/handlers"
	"github.com/cyberlab/backend/internal/middleware"
	"github.com/cyberlab/backend/pkg/ws"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all routes
func SetupRouter(
	jwtSecret string,
	authHandler *handlers.AuthHandler,
	challengeHandler *handlers.ChallengeHandler,
	containerHandler *handlers.ContainerHandler,
	submissionHandler *handlers.SubmissionHandler,
	rankingHandler *handlers.RankingHandler,
	announcementHandler *handlers.AnnouncementHandler,
	aiHandler *handlers.AIHandler,
	adminHandler *handlers.AdminHandler,
	wsHub *ws.Hub,
) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket
	r.GET("/ws", func(c *gin.Context) {
		userID := c.Query("userId")
		if userID == "" {
			c.JSON(401, gin.H{"error": "missing userId"})
			return
		}
		var uid uint
		if _, err := fmt.Sscanf(userID, "%d", &uid); err != nil {
			c.JSON(401, gin.H{"error": "invalid userId"})
			return
		}
		wsHub.HandleWebSocket(c.Writer, c.Request, uid)
	})

	api := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Challenges (public)
		challenges := api.Group("/challenges")
		{
			challenges.GET("", challengeHandler.GetChallenges)
			challenges.GET("/search", challengeHandler.SearchChallenges)
			challenges.GET("/categories", challengeHandler.GetCategories)
			challenges.GET("/:id", challengeHandler.GetChallenge)
		}

		// Announcements (public)
		announcements := api.Group("/announcement")
		{
			announcements.GET("", announcementHandler.GetAnnouncements)
			announcements.GET("/:id", announcementHandler.GetAnnouncement)
		}

		// Ranking (public)
		ranking := api.Group("/ranking")
		{
			ranking.GET("/global", rankingHandler.GetGlobalRanking)
			ranking.GET("/week", rankingHandler.GetWeeklyRanking)
			ranking.GET("/month", rankingHandler.GetMonthlyRanking)
		}

		// Authenticated routes
		authed := api.Group("")
		authed.Use(middleware.JWTAuth(jwtSecret))
		{
			// Profile
			authed.GET("/auth/profile", authHandler.GetProfile)
			authed.PUT("/auth/profile", authHandler.UpdateProfile)

			// Container
			authed.POST("/container/start", containerHandler.StartChallenge)
			authed.GET("/container/status/:id", containerHandler.GetContainerStatus)
			authed.POST("/container/stop", containerHandler.StopChallenge)
			authed.GET("/container/instances", containerHandler.GetUserInstances)

			// Submit
			authed.POST("/submit", submissionHandler.SubmitFlag)
			authed.GET("/submit/history", submissionHandler.GetSubmissionHistory)
			authed.GET("/submit/solved", submissionHandler.GetSolvedChallenges)

			// AI
			authed.POST("/ai/hint", aiHandler.GetHint)
			authed.POST("/ai/audit", aiHandler.AuditCode)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(jwtSecret), middleware.AdminOnly())
		{
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/logs", adminHandler.GetLogs)

			// Challenge management
			admin.POST("/challenge", challengeHandler.CreateChallenge)
			admin.PUT("/challenge/:id", challengeHandler.UpdateChallenge)
			admin.DELETE("/challenge/:id", challengeHandler.DeleteChallenge)

			// Container management
			admin.GET("/container", adminHandler.GetContainers)
			admin.DELETE("/container/:id", adminHandler.ForceDestroyContainer)

			// Announcement management
			admin.POST("/announcement", adminHandler.CreateAnnouncement)
			admin.DELETE("/announcement/:id", adminHandler.DeleteAnnouncement)
		}
	}

	return r
}
