package router

import (
	"github.com/cyberlab/backend/internal/handlers"
	"github.com/cyberlab/backend/internal/middleware"
	"github.com/cyberlab/backend/pkg/metrics"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/cyberlab/backend/pkg/ws"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes with security hardening.
func SetupRouter(
	jwtSecret string,
	serverMode string,
	allowedOrigins []string,
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
	// Set Gin mode
	gin.SetMode(serverMode)

	r := gin.New()
	r.Use(gin.Recovery())

	// Security headers middleware
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	})

	// CORS middleware — restricted to specific origins
	originMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		originMap[o] = true
	}

	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		if len(originMap) == 0 || originMap["*"] {
			allowOrigin = "*"
		} else if originMap[origin] {
			allowOrigin = origin
		}

		if allowOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowOrigin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

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

	// Prometheus metrics
	r.GET("/metrics", metrics.Handler())

	// Global metrics middleware
	r.Use(metrics.Middleware())

	// Swagger — only available in debug/test mode
	if serverMode == "debug" {
		// Register swagger endpoints here if needed
		utils.Log.Warn("Swagger UI is accessible — disable in production")
	}

	// WebSocket — JWT authenticated via query parameter
	r.GET("/ws", func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			c.JSON(401, gin.H{"error": "missing token"})
			return
		}

		claims, err := utils.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		wsHub.HandleWebSocket(c.Writer, c.Request, claims.UserID)
	})

	api := r.Group("/api/v1")
	{
		// Auth routes (public, rate-limited)
		auth := api.Group("/auth")
		{
			auth.POST("/register", middleware.RateLimitMiddleware(middleware.RegisterLimiter), authHandler.Register)
			auth.POST("/login", middleware.RateLimitMiddleware(middleware.LoginLimiter), authHandler.Login)
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
