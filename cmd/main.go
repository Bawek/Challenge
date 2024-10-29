package main

import (
	"log/slog"
	"os"

	_ "github.com/izymalhaw/go-crud/yishakterefe/docs"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/api/handlers"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/config"
	customlogger "github.com/izymalhaw/go-crud/yishakterefe/internal/core/logger"
	"github.com/izymalhaw/go-crud/yishakterefe/internal/repository"
	person_service "github.com/izymalhaw/go-crud/yishakterefe/internal/services/person"
)

const (
	version = "1.0.0"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		os.Exit(1)
	}
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	logger := customlogger.NewLogger(cfg.Env, cfg.LogLevel, version)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	store := repository.NewInMemoryUserRepo()
	personService := person_service.NewPersonSvc(store)
	webSrv := handlers.NewApp(cfg.Port, personService, logger)
	logger.Info("server running ")
	webSrv.Run()
}
