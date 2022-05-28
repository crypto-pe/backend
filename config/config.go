package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	GitCommit string `toml:"-"`
	Mode      Mode   `toml:"-"`

	Service ServiceConfig `toml:"service"`
	Logging LoggingConfig `toml:"logging"`
	Auth    Auth          `toml:"auth"`

	DB DBConfig `toml:"db"`
}

type ServiceConfig struct {
	// Name that identifies the application
	Name string `toml:"name"`

	// Listen network url for the HTTP/RPC server
	Listen string `toml:"listen"`

	// Mode is the operating mode of the application, one of:
	// "development", "dev", "production" or "prod"
	Mode string `toml:"mode"`
}

type LoggingConfig struct {
	Level   string `toml:"level"`
	JSON    bool   `toml:"json"`
	Concise bool   `toml:"concise"`
}

type DBConfig struct {
	Host     string `toml:"host"`
	Database string `toml:"database"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Port     string `toml:"port"`
}

func (cfg *Config) DBString() string {
	if cfg.Mode == DevelopmentMode {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Host, cfg.DB.Database)
}

type Auth struct {
	JWTSecret string `toml:"jwt_secret"`
}

func NewFromFile(file string, env string, cfg *Config) error {
	if file == "" {
		file = env
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return fmt.Errorf("failed to load config file: %w", err)
	}
	if _, err := toml.DecodeFile(file, cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return InitConfig(cfg)
}

func InitConfig(cfg *Config) error {
	// Service mode
	var mode Mode
	switch cfg.Service.Mode {
	case "dev", "development":
		mode = DevelopmentMode
	case "prod", "production":
		mode = ProductionMode
	default:
		return fmt.Errorf("config service.mode value is invalid, must be one of \"development\", \"dev\", \"production\" or \"prod\"")
	}
	cfg.Mode = mode
	cfg.Service.Mode = mode.String()

	// Validate auth
	if cfg.Auth.JWTSecret == "" || len(cfg.Auth.JWTSecret) < 10 {
		return fmt.Errorf("config auth.jwt_secret must be at least 10 characters long")
	}

	return nil
}

type Mode uint32

const (
	DevelopmentMode Mode = iota
	ProductionMode
)

func (m Mode) String() string {
	switch m {
	case DevelopmentMode:
		return "development"
	case ProductionMode:
		return "production"
	default:
		return ""
	}
}
