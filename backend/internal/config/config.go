package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Docker   DockerConfig   `yaml:"docker"`
	Log      LogConfig      `yaml:"log"`
	AI       AIConfig       `yaml:"ai"`
	Security SecurityConfig `yaml:"security"`
}

type SecurityConfig struct {
	FlagHashSalt string `yaml:"flag_hash_salt"`
}

type ServerConfig struct {
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type DockerConfig struct {
	Host                string   `yaml:"host"`
	APIVersion          string   `yaml:"api_version"`
	PortRangeStart      int      `yaml:"port_range_start"`
	PortRangeEnd        int      `yaml:"port_range_end"`
	DefaultTimeoutMin   int      `yaml:"default_timeout_minutes"`
	CPULimit            float64  `yaml:"cpu_limit"`
	MemoryLimitMB       int64    `yaml:"memory_limit_mb"`
	NetworkName         string   `yaml:"network_name"`
	AllowedImages       []string `yaml:"allowed_images"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

type AIConfig struct {
	Provider    string  `yaml:"provider"`
	APIKey      string  `yaml:"api_key"`
	Model       string  `yaml:"model"`
	MaxTokens   int     `yaml:"max_tokens"`
	Temperature float64 `yaml:"temperature"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override secrets from environment variables (highest priority)
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("AI_API_KEY"); v != "" {
		cfg.AI.APIKey = v
	}
	if v := os.Getenv("FLAG_HASH_SALT"); v != "" {
		cfg.Security.FlagHashSalt = v
	}

	// Validate required secrets
	if cfg.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET is required (set in config.yaml or env var)")
	}
	if cfg.Database.Password == "" {
		return nil, errors.New("DB_PASSWORD is required (set in config.yaml or env var)")
	}

	// Default allowed images
	if len(cfg.Docker.AllowedImages) == 0 {
		cfg.Docker.AllowedImages = []string{"cyberlab/*"}
	}

	return &cfg, nil
}
