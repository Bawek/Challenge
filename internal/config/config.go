package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	customlogger "github.com/izymalhaw/go-crud/yishakterefe/internal/core/logger"
)

var (
	ErrInvalidPort  = errors.New("port number is invalid")
	ErrLogLevel     = errors.New("log level not set")
	ErrInvalidEnv   = errors.New("env not set or invalid")
	ErrInvalidLevel = errors.New("invalid log level")
)

type Config struct {
	Port          int
	LogLevel      slog.Level
	Env           string
}

var Environment = map[string]string{
	"dev":  "development",
	"prod": "production",
}

func (c *Config) loadEnv() error {
	env := os.Getenv("ENV")
	if env == "" {
		return ErrInvalidEnv
	}

	evalue, ok := Environment[env]
	if !ok {
		return ErrInvalidEnv
	}
	c.Env = evalue

	
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return ErrLogLevel
	}

	lvl, ok := customlogger.LogLevels[logLevel]
	if !ok {
		return ErrInvalidLevel
	}
	c.LogLevel = lvl

	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ErrInvalidPort
	}
	c.Port = port
	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := config.loadEnv(); err != nil {
		return nil, err
	}
	return config, nil
}
