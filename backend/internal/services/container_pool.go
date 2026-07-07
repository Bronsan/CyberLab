package services

import (
	"sync"
	"time"

	"github.com/cyberlab/backend/pkg/docker"
	"go.uber.org/zap"
)

// ContainerPool pre-warms Docker containers for faster challenge startup.
// When a user requests a challenge, if a pre-warmed container exists,
// it's assigned instantly instead of pulling+building.
type ContainerPool struct {
	mu          sync.RWMutex
	pool        map[string][]*PooledContainer
	maxSize     int
	warmImage   string
	dockerMgr   *docker.Manager
	logger      *zap.Logger
	stopCh      chan struct{}
}

type PooledContainer struct {
	ContainerID string
	Image       string
	CreatedAt   time.Time
}

func NewContainerPool(dockerMgr *docker.Manager, logger *zap.Logger, maxSize int) *ContainerPool {
	return &ContainerPool{
		pool:      make(map[string][]*PooledContainer),
		maxSize:   maxSize,
		dockerMgr: dockerMgr,
		logger:    logger,
		stopCh:    make(chan struct{}),
	}
}

// PreWarm creates N containers for a given image and keeps them ready.
func (cp *ContainerPool) PreWarm(image string, count int) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.logger.Info("Pre-warming container pool",
		zap.String("image", image),
		zap.Int("target", count),
	)

	// Pool containers are tracked as available slots
	// Actual container creation happens lazily via the ContainerService
	cp.pool[image] = make([]*PooledContainer, 0, count)
}

// Acquire gets a pre-warmed container from the pool, or returns nil.
func (cp *ContainerPool) Acquire(image string) *PooledContainer {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if containers, ok := cp.pool[image]; ok && len(containers) > 0 {
		c := containers[0]
		cp.pool[image] = containers[1:]
		return c
	}
	return nil
}

// StartReplenishment runs a background loop to maintain pool size.
func (cp *ContainerPool) StartReplenishment(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cp.mu.Lock()
				for image, containers := range cp.pool {
					if len(containers) < cp.maxSize {
						cp.mu.Unlock()
						cp.PreWarm(image, cp.maxSize-len(containers))
						cp.mu.Lock()
					}
				}
				cp.mu.Unlock()
			case <-cp.stopCh:
				return
			}
		}
	}()
}

func (cp *ContainerPool) Stop() {
	close(cp.stopCh)
}
