package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"users/config"
	"users/initializers/postgre"
	"users/pkg/zaplogger"
	"users/service"
	"users/storage/dbpostgre"
	"users/transport/http"
)

func main() {
	configFile := "conf_local"
	if os.Getenv("ENV") == "docker" {
		configFile = "conf"
	}

	// Viper
	_, cfg, errViper := config.NewViper(configFile)
	if errViper != nil {
		log.Fatal(errors.WithMessage(errViper, "Viper startup error"))
	}

	// Zap logger
	logger, loggerCleanup, errZapLogger := zaplogger.New(zaplogger.Mode(cfg.Logger.Development))
	if errZapLogger != nil {
		log.Fatal(errors.WithMessage(errZapLogger, "Zap logger startup error"))
	}

	// Postgre
	db, postgCleanup, err := postgre.NewDB(cfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to DB", zap.Error(err))
	}

	dataBaseRepo := dbpostgre.NewDataBaseRepositoryImpl(db)
	dataBaseWorker := service.NewDataBaseWorker(dataBaseRepo)

	userController := http.NewUserController(logger, dataBaseWorker)

	router := http.NewRouter(userController, logger, cfg)
	router.RegisterRoutes()

	// Channel for error transmission
	errCh := make(chan error, 1)

	// Router in goroutine
	go func() {
		err := router.Start()
		if err != nil {
			logger.Error("Error starting router", zap.Error(err))
		}
		errCh <- err
	}()

	// Handle shutdown gracefully
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		logger.Error(err.Error())
	case <-shutdown:
	}
	logger.Info("Received shutdown signal. Shutting down...")
	loggerCleanup()
	postgCleanup()
	logger.Info("Application stopped gracefully")

}
