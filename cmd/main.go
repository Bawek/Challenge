package main

import (
	"log/slog"
	"os"

	// "github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/izymalhaw/go-crud/yishakterefe/docs"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/api/handlers"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/config"
	customlogger "github.com/izymalhaw/go-crud/yishakterefe/internal/core/logger"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/repository"
	person_service "github.com/izymalhaw/go-crud/yishakterefe/internal/services/person"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const version = "1.0.0"

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
	}

	// Load configuration

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Connect to MySQL
	db, err := gorm.Open(mysql.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to Load DB ", "error", err)
	}
	defer sqlDB.Close()

	// Logger
	logger := customlogger.NewLogger(cfg.Env, cfg.LogLevel, version)

	// Repository & Service
	// store := repository.NewPostgresPersonRepo(db) // pass db
	store := repository.NewMySqlPersonRepo(db) // pass db
	personService := person_service.NewPersonSvc(store)

	// Start server
	webSrv := handlers.NewApp(cfg.Port, personService, logger)
	logger.Info("server running", "port", cfg.Port)
	webSrv.Run()
}
