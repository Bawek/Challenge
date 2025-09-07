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
	ErrDBUrl        = errors.New("invalid database url")
)

type Config struct {
	Port     int
	LogLevel slog.Level
	Env      string
	DBUrl    string
}

var Environment = map[string]string{
	"dev":  "development",
	"prod": "production",
}

func (c *Config) loadEnv() error {
	// Load ENV
	env := os.Getenv("ENV")
	if env == "" {
		return ErrInvalidEnv
	}
	evalue, ok := Environment[env]
	if !ok {
		return ErrInvalidEnv
	}
	c.Env = evalue

	// Load LOG_LEVEL
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return ErrLogLevel
	}
	lvl, ok := customlogger.LogLevels[logLevel]
	if !ok {
		return ErrInvalidLevel
	}
	c.LogLevel = lvl

	// Load PORT
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ErrInvalidPort
	}
	c.Port = port

	// Load DATABASE_URL
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return ErrDBUrl
	}
	c.DBUrl = dburl

	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := config.loadEnv(); err != nil {
		return nil, err
	}
	return config, nil
}
