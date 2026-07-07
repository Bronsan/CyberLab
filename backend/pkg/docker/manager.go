package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Manager wraps Docker CLI operations
type Manager struct {
	networkName    string
	cpuLimit       float64
	memoryLimitMB  int64
	portRangeStart int
	portRangeEnd   int
}

type Config struct {
	Host           string
	APIVersion     string
	NetworkName    string
	CPULimit       float64
	MemoryLimitMB  int64
	PortRangeStart int
	PortRangeEnd   int
}

func NewManager(cfg *Config) (*Manager, error) {
	// Test Docker is available
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info", "--format", "{{.ServerVersion}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker not available: %w", err)
	}

	_ = output // version info

	return &Manager{
		networkName:    cfg.NetworkName,
		cpuLimit:       cfg.CPULimit,
		memoryLimitMB:  cfg.MemoryLimitMB,
		portRangeStart: cfg.PortRangeStart,
		portRangeEnd:   cfg.PortRangeEnd,
	}, nil
}

// CreateAndStartContainer creates and starts a container from the given image
func (m *Manager) CreateAndStartContainer(ctx context.Context, imageName string, containerName string, hostPort int, containerPort string) (string, error) {
	// Pull image first
	pullCmd := exec.CommandContext(ctx, "docker", "pull", imageName)
	if output, err := pullCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to pull image %s: %s: %w", imageName, string(output), err)
	}

	// Create and start container
	args := []string{
		"run", "-d",
		"--name", containerName,
		"--network", m.networkName,
		"--cpus", fmt.Sprintf("%f", m.cpuLimit),
		"--memory", fmt.Sprintf("%dm", m.memoryLimitMB),
		"-p", fmt.Sprintf("%d:%s", hostPort, containerPort),
		"--restart", "no",
		imageName,
	}

	runCmd := exec.CommandContext(ctx, "docker", args...)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to create/start container: %s: %w", string(output), err)
	}

	containerID := strings.TrimSpace(string(output))
	return containerID, nil
}

// StopContainer stops a container
func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	cmd := exec.CommandContext(ctx, "docker", "stop", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop container %s: %s: %w", containerID, string(output), err)
	}
	return nil
}

// RemoveContainer removes a container
func (m *Manager) RemoveContainer(ctx context.Context, containerID string) error {
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", containerID)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove container %s: %s: %w", containerID, string(output), err)
	}
	return nil
}

// GetContainerStatus returns the container's status
func (m *Manager) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.State.Status}}", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to inspect container %s: %s: %w", containerID, string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// ContainerExists checks if a container exists
func (m *Manager) ContainerExists(ctx context.Context, containerID string) bool {
	cmd := exec.CommandContext(ctx, "docker", "inspect", containerID)
	return cmd.Run() == nil
}

// GetContainerLogs returns container logs
func (m *Manager) GetContainerLogs(ctx context.Context, containerID string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "logs", "--tail", "100", containerID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs for container %s: %s: %w", containerID, string(output), err)
	}
	return string(output), nil
}

// ListContainers lists all running containers for this project
func (m *Manager) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	cmd := exec.CommandContext(ctx, "docker", "ps", "--format", `{"id":"{{.ID}}","name":"{{.Names}}","image":"{{.Image}}","status":"{{.Status}}","ports":"{{.Ports}}"}`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %s: %w", string(output), err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containers []ContainerInfo
	for _, line := range lines {
		if line == "" {
			continue
		}
		var info ContainerInfo
		if err := json.Unmarshal([]byte(line), &info); err != nil {
			continue
		}
		containers = append(containers, info)
	}
	return containers, nil
}

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
	Ports  string `json:"ports"`
}

// Ping checks if Docker daemon is accessible
func (m *Manager) Ping(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "docker", "info")
	_, err := cmd.CombinedOutput()
	return err
}

// Close is a no-op for CLI-based manager
func (m *Manager) Close() error {
	return nil
}
