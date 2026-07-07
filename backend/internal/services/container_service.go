package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cyberlab/backend/internal/models"
	"github.com/cyberlab/backend/internal/repository"
	"github.com/cyberlab/backend/pkg/docker"
	"github.com/cyberlab/backend/pkg/utils"
	ws "github.com/cyberlab/backend/pkg/ws"
	"go.uber.org/zap"
)

type ContainerService struct {
	instanceRepo  *repository.InstanceRepository
	challengeRepo *repository.ChallengeRepository
	userRepo      *repository.UserRepository
	logRepo       *repository.LogRepository
	dockerMgr     *docker.Manager
	hub           *ws.Hub
	logger        *zap.Logger
	portPool      map[int]bool
	portMu        sync.Mutex
	portStart     int
	portEnd       int
	allowedImages []string
}

func NewContainerService(
	instanceRepo *repository.InstanceRepository,
	challengeRepo *repository.ChallengeRepository,
	userRepo *repository.UserRepository,
	logRepo *repository.LogRepository,
	dockerMgr *docker.Manager,
	hub *ws.Hub,
	logger *zap.Logger,
	portStart, portEnd int,
	allowedImages []string,
) *ContainerService {
	// Initialize port pool
	pool := make(map[int]bool)
	for p := portStart; p <= portEnd; p++ {
		pool[p] = true
	}

	return &ContainerService{
		instanceRepo:  instanceRepo,
		challengeRepo: challengeRepo,
		userRepo:      userRepo,
		logRepo:       logRepo,
		dockerMgr:     dockerMgr,
		hub:           hub,
		logger:        logger,
		portPool:      pool,
		portStart:     portStart,
		portEnd:       portEnd,
		allowedImages: allowedImages,
	}
}

type StartContainerResponse struct {
	InstanceID uint   `json:"instanceId"`
	URL        string `json:"url"`
}

func (s *ContainerService) StartContainer(userID uint, challengeID uint) (*StartContainerResponse, error) {
	// Check if user already has a running instance for this challenge
	existing, err := s.instanceRepo.FindByUserAndChallenge(userID, challengeID)
	if err == nil && existing != nil {
		return nil, errors.New("you already have a running instance for this challenge")
	}

	// Get challenge info
	challenge, err := s.challengeRepo.FindByID(challengeID)
	if err != nil {
		return nil, errors.New("challenge not found")
	}

	if !challenge.IsActive {
		return nil, errors.New("challenge is not active")
	}

	if challenge.DockerImage == "" {
		return nil, errors.New("challenge has no docker image configured")
	}

	// Validate image against allowed list
	if err := s.validateAllowedImage(challenge.DockerImage); err != nil {
		return nil, err
	}

	// Allocate port
	port, err := s.allocatePort()
	if err != nil {
		return nil, errors.New("no available ports")
	}

	containerName := fmt.Sprintf("cyberlab-%d-%d-%d", userID, challengeID, time.Now().Unix())

	// Create and start container via Docker
	ctx := context.Background()
	containerID, err := s.dockerMgr.CreateAndStartContainer(
		ctx,
		challenge.DockerImage,
		containerName,
		port,
		"80", // Default container port
	)
	if err != nil {
		s.releasePort(port)
		s.logger.Error("Failed to start container", zap.Error(err))
		return nil, errors.New("failed to start container: " + err.Error())
	}

	// Create instance record
	now := time.Now()
	expireTime := now.Add(time.Duration(challenge.TimeoutMinutes) * time.Minute)

	instance := &models.ChallengeInstance{
		UserID:        userID,
		ChallengeID:   challengeID,
		ContainerID:   containerID,
		ContainerName: containerName,
		HostPort:      port,
		Status:        models.InstanceStatusRunning,
		StartTime:     now,
		ExpireTime:    expireTime,
	}

	if err := s.instanceRepo.Create(instance); err != nil {
		// Rollback: stop and remove container
		s.dockerMgr.StopContainer(ctx, containerID)
		s.dockerMgr.RemoveContainer(ctx, containerID)
		s.releasePort(port)
		s.logger.Error("Failed to save instance", zap.Error(err))
		return nil, errors.New("failed to save instance")
	}

	// Log
	s.logRepo.Create(&models.SystemLog{
		UserID:    userID,
		Action:    models.ActionStartContainer,
		IP:        "",
	})

	// WebSocket notification
	s.hub.SendToUser(userID, ws.Message{
		Event: "container_created",
		Data: map[string]interface{}{
			"instanceId": instance.ID,
			"status":     models.InstanceStatusRunning,
		},
	})

	url := fmt.Sprintf("http://lab.cyberlab.com:%d", port)
	return &StartContainerResponse{
		InstanceID: instance.ID,
		URL:        url,
	}, nil
}

func (s *ContainerService) GetContainerStatus(instanceID uint) (*models.ChallengeInstance, error) {
	instance, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return nil, errors.New("instance not found")
	}

	// Update status from Docker
	ctx := context.Background()
	status, err := s.dockerMgr.GetContainerStatus(ctx, instance.ContainerID)
	if err == nil {
		dockerStatus := instance.Status
		switch status {
		case "running":
			dockerStatus = models.InstanceStatusRunning
		case "exited", "stopped":
			dockerStatus = models.InstanceStatusStopped
		case "created":
			dockerStatus = models.InstanceStatusCreated
		}

		if dockerStatus != instance.Status {
			s.instanceRepo.UpdateStatus(instance.ID, dockerStatus)
			instance.Status = dockerStatus
		}
	}

	return instance, nil
}

func (s *ContainerService) StopContainer(userID uint, instanceID uint) error {
	instance, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return errors.New("instance not found")
	}

	if instance.UserID != userID {
		return errors.New("instance does not belong to this user")
	}

	if instance.Status != models.InstanceStatusRunning {
		return errors.New("instance is not running")
	}

	ctx := context.Background()
	if err := s.dockerMgr.StopContainer(ctx, instance.ContainerID); err != nil {
		s.logger.Error("Failed to stop container", zap.Error(err))
		return errors.New("failed to stop container")
	}

	s.instanceRepo.UpdateStatus(instance.ID, models.InstanceStatusStopped)

	// WebSocket notification
	s.hub.SendToUser(userID, ws.Message{
		Event: "container_stopped",
		Data: map[string]interface{}{
			"instanceId": instanceID,
			"status":     models.InstanceStatusStopped,
		},
	})

	s.logRepo.Create(&models.SystemLog{
		UserID: userID,
		Action: models.ActionStopContainer,
	})

	return nil
}

func (s *ContainerService) DestroyContainer(instanceID uint) error {
	instance, err := s.instanceRepo.FindByID(instanceID)
	if err != nil {
		return errors.New("instance not found")
	}

	ctx := context.Background()

	// Stop and remove container
	s.dockerMgr.StopContainer(ctx, instance.ContainerID)
	s.dockerMgr.RemoveContainer(ctx, instance.ContainerID)

	// Update status
	s.instanceRepo.UpdateStatus(instance.ID, models.InstanceStatusDestroyed)

	// Release port
	s.releasePort(instance.HostPort)

	// WebSocket notification
	s.hub.SendToUser(instance.UserID, ws.Message{
		Event: "container_destroyed",
		Data: map[string]interface{}{
			"instanceId": instanceID,
		},
	})

	return nil
}

func (s *ContainerService) RecycleExpiredContainers() {
	instances, err := s.instanceRepo.FindExpired()
	if err != nil {
		s.logger.Error("Failed to find expired instances", zap.Error(err))
		return
	}

	for _, instance := range instances {
		s.logger.Info("Recycling expired container",
			zap.Uint("instanceId", instance.ID),
			zap.String("containerId", instance.ContainerID),
		)
		if err := s.DestroyContainer(instance.ID); err != nil {
			s.logger.Error("Failed to recycle container", zap.Error(err))
		}
	}
}

func (s *ContainerService) GetUserInstances(userID uint) ([]models.ChallengeInstance, error) {
	return s.instanceRepo.FindRunningByUser(userID)
}

func (s *ContainerService) GetAllInstances(page, pageSize int) ([]models.ChallengeInstance, int64, error) {
	return s.instanceRepo.FindAll(page, pageSize)
}

func (s *ContainerService) GetContainerStats() (map[string]int64, error) {
	return s.instanceRepo.CountByStatus()
}

func (s *ContainerService) allPortsInUse() bool {
	s.portMu.Lock()
	defer s.portMu.Unlock()
	return len(s.portPool) == 0
}

func (s *ContainerService) allocatePort() (int, error) {
	s.portMu.Lock()
	defer s.portMu.Unlock()

	if len(s.portPool) == 0 {
		return 0, errors.New("port pool exhausted")
	}

	// Random allocation
	var ports []int
	for p := range s.portPool {
		ports = append(ports, p)
	}

	idx := rand.Intn(len(ports))
	port := ports[idx]
	delete(s.portPool, port)

	return port, nil
}

func (s *ContainerService) releasePort(port int) {
	s.portMu.Lock()
	defer s.portMu.Unlock()
	s.portPool[port] = true
}

// validateAllowedImage checks that the image matches the allowed list.
// Supports glob patterns like "cyberlab/*" or exact matches like "library/nginx:latest".
func (s *ContainerService) validateAllowedImage(image string) error {
	if len(s.allowedImages) == 0 {
		return nil // no restrictions
	}

	for _, pattern := range s.allowedImages {
		matched, err := matchGlob(pattern, image)
		if err == nil && matched {
			return nil
		}
	}

	utils.Log.Warn("Blocked attempt to use unauthorized Docker image",
		zap.String("image", image),
		zap.Strings("allowed", s.allowedImages),
	)
	return fmt.Errorf("docker image %s is not in the allowed list", image)
}

// Simple glob match (supports "*" prefix/suffix matching)
func matchGlob(pattern, str string) (bool, error) {
	// Exact match
	if pattern == str {
		return true, nil
	}

	// Wildcard: "prefix/*" matches "prefix/anything"
	parts := patternSplit(pattern)
	strParts := patternSplit(str)

	if len(parts) == 2 && parts[0] == strParts[0] && parts[1] == "*" {
		return true, nil
	}
	if len(parts) == 2 && parts[1] == "*" && parts[0] == strParts[0] {
		return true, nil
	}

	return false, nil
}

func patternSplit(s string) []string {
	var parts []string
	current := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	parts = append(parts, current)
	return parts
}
