package scheduler

import (
	"time"

	"github.com/cyberlab/backend/internal/services"
	"go.uber.org/zap"
)

type RecycleWorker struct {
	containerService *services.ContainerService
	logger           *zap.Logger
	interval         time.Duration
	stopChan         chan struct{}
}

func NewRecycleWorker(containerService *services.ContainerService, logger *zap.Logger) *RecycleWorker {
	return &RecycleWorker{
		containerService: containerService,
		logger:           logger,
		interval:         1 * time.Minute,
		stopChan:         make(chan struct{}),
	}
}

// Start begins the periodic recycling loop
func (w *RecycleWorker) Start() {
	w.logger.Info("Container recycle worker started", zap.Duration("interval", w.interval))

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.logger.Debug("Running container recycle check")
				w.containerService.RecycleExpiredContainers()
			case <-w.stopChan:
				w.logger.Info("Container recycle worker stopped")
				return
			}
		}
	}()
}

// Stop gracefully stops the recycle worker
func (w *RecycleWorker) Stop() {
	close(w.stopChan)
}
