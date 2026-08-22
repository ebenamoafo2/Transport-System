package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port            int           `yaml:"port"`
	ReadTimeoutSec  int           `yaml:"read_timeout_sec"`
	WriteTimeoutSec int           `yaml:"write_timeout_sec"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	ConnMaxLifetime     int `yaml:"conn_max_lifetime_sec"`
	ConnMaxIdleTime     int `yaml:"conn_max_idle_time_sec"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig reads and parses the configuration file from the given path.
// It returns a Config struct or an error if it fails.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	//Override the server port with SERVER_PORT if provided
	if port := os.Getenv("SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
		}
		cfg.Server.Port = p
	}

	// Override the database password with DB_PASSWORD if provided.
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	//Validate the configs
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &cfg, nil
}

// Validate make sure the valid configs are used
func (c Config) Validate() error {
	var errs []error
	if err := c.validateServer(); err != nil {
		errs = append(errs, fmt.Errorf("server: %w", err))
	}
	if err := c.validateDatabase(); err != nil {
		errs = append(errs, fmt.Errorf("database: %w", err))
	}
	return errors.Join(errs...)
}

// validateServer validates the HTTP server configuration.
func (c Config) validateServer() error {
	var errs []error

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		errs = append(errs, fmt.Errorf("port must be between 1 and 65535"))
	}

	if c.Server.ReadTimeoutSec <= 0 {
		errs = append(errs, fmt.Errorf("read_timeout_sec must be greater than 0"))
	}

	if c.Server.WriteTimeoutSec <= 0 {
		errs = append(errs, fmt.Errorf("write_timeout_sec must be greater than 0"))
	}
	if c.Server.IdleTimeout <= 0 {
		errs = append(errs, fmt.Errorf("idle_timeout_sec must be greated than 0"))
	}

	return errors.Join(errs...)
}

// validateDatabase validates the database configuration.
func (c Config) validateDatabase() error {
	var errs []error

	if c.Database.Host == "" {
		errs = append(errs, fmt.Errorf("host cannot be empty"))
	}

	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		errs = append(errs, fmt.Errorf("port must be between 1 and 65535"))
	}

	if c.Database.User == "" {
		errs = append(errs, fmt.Errorf("user cannot be empty"))
	}

	if c.Database.Password == "" {
		errs = append(errs, fmt.Errorf("password cannot be empty"))
	}

	if c.Database.Name == "" {
		errs = append(errs, fmt.Errorf("name cannot be empty"))
	}

	if c.Database.MaxOpenConns <= 0 {
		errs = append(errs, fmt.Errorf("max_open_conns must be greater than 0"))
	}

	if c.Database.MaxIdleConns < 0 {
		errs = append(errs, fmt.Errorf("max_idle_conns cannot be negative"))
	}
	// You cannot have more idle connections than the total
	// number of connections allowed.
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		errs = append(errs, fmt.Errorf(
			"max_idle_conns cannot be greater than max_open_conns",
		))
	}

	return errors.Join(errs...)
}
