package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberlab/backend/internal/config"
	"github.com/cyberlab/backend/internal/handlers"
	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/internal/router"
	"github.com/cyberlab/backend/internal/scheduler"
	"github.com/cyberlab/backend/internal/services"
	"github.com/cyberlab/backend/pkg/docker"
	"github.com/cyberlab/backend/pkg/utils"
	"github.com/cyberlab/backend/pkg/ws"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// @title CyberLab API
// @version 1.0
// @description CyberLab 网络安全靶场平台 API
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfgPath := "config/config.yaml"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		cfgPath = v
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	utils.InitLogger(cfg.Log.Level, cfg.Log.Filename, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)
	utils.Log.Info("CyberLab starting...", zap.String("mode", cfg.Server.Mode))

	// Connect to MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		utils.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Apply connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		utils.Log.Fatal("Failed to get SQL DB from GORM", zap.Error(err))
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	utils.Log.Info("Database connected", zap.String("host", cfg.Database.Host))

	// Auto migrate
	db.AutoMigrate(
		&models.User{},
		&models.Challenge{},
		&models.ChallengeTag{},
		&models.ChallengeInstance{},
		&models.Submission{},
		&models.UserChallenge{},
		&models.Announcement{},
		&models.AIHint{},
		&models.SystemLog{},
	)

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		utils.Log.Warn("Redis connection failed, caching disabled", zap.Error(err))
	} else {
		utils.Log.Info("Redis connected")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	challengeRepo := repository.NewChallengeRepository(db)
	instanceRepo := repository.NewInstanceRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	logRepo := repository.NewLogRepository(db)
	aiRepo := repository.NewAIRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)

	// Initialize Docker manager
	dockerMgr, err := docker.NewManager(&docker.Config{
		Host:           cfg.Docker.Host,
		APIVersion:     cfg.Docker.APIVersion,
		NetworkName:    cfg.Docker.NetworkName,
		CPULimit:       cfg.Docker.CPULimit,
		MemoryLimitMB:  cfg.Docker.MemoryLimitMB,
		PortRangeStart: cfg.Docker.PortRangeStart,
		PortRangeEnd:   cfg.Docker.PortRangeEnd,
	})
	if err != nil {
		utils.Log.Warn("Docker initialization failed, container features disabled", zap.Error(err))
		dockerMgr = nil
	} else {
		utils.Log.Info("Docker manager initialized")
	}

	// Initialize WebSocket hub
	wsHub := ws.NewHub(utils.Log)
	go wsHub.Run()

	// Initialize services
	authService := services.NewAuthService(userRepo, logRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	challengeService := services.NewChallengeService(challengeRepo)
	containerService := services.NewContainerService(
		instanceRepo, challengeRepo, userRepo, logRepo,
		dockerMgr, wsHub, utils.Log,
		cfg.Docker.PortRangeStart, cfg.Docker.PortRangeEnd,
		cfg.Docker.AllowedImages,
	)
	rankingService := services.NewRankingService(userRepo, redisClient)
	aiService := services.NewAIService(
		services.AIServiceConfig{
			Provider:    cfg.AI.Provider,
			APIKey:      cfg.AI.APIKey,
			Model:       cfg.AI.Model,
			MaxTokens:   cfg.AI.MaxTokens,
			Temperature: cfg.AI.Temperature,
		},
		aiRepo, logRepo,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	challengeHandler := handlers.NewChallengeHandler(challengeService, cfg.Security.FlagHashSalt)
	containerHandler := handlers.NewContainerHandler(containerService)
	submissionHandler := handlers.NewSubmissionHandler(submissionRepo, challengeRepo, rankingService, logRepo, instanceRepo, cfg.Security.FlagHashSalt)
	rankingHandler := handlers.NewRankingHandler(rankingService)
	announcementHandler := handlers.NewAnnouncementHandler(announcementRepo)
	aiHandler := handlers.NewAIHandler(aiService)
	adminHandler := handlers.NewAdminHandler(userRepo, containerService, challengeService, announcementRepo, logRepo)

	// Setup router
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"https://bronsan.github.io",
	}
	r := router.SetupRouter(
		cfg.JWT.Secret,
		cfg.Server.Mode,
		allowedOrigins,
		authHandler, challengeHandler, containerHandler,
		submissionHandler, rankingHandler, announcementHandler,
		aiHandler, adminHandler, wsHub,
	)

	// Start recycle worker
	recycleWorker := scheduler.NewRecycleWorker(containerService, utils.Log)
	recycleWorker.Start()

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()
	utils.Log.Info("Server started", zap.String("addr", addr))

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Log.Info("Shutting down server...")
	recycleWorker.Stop()
	if dockerMgr != nil {
		dockerMgr.Close()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		utils.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	utils.Log.Info("Server exited")
}
