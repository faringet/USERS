package postgre

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"users/config"
	"users/domain"
)

func NewDB(cfg *config.Config, logger *zap.Logger) (db *gorm.DB, cleanup func() error, err error) {
	dsn := cfg.Postgre.DbURL

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to DB:", zap.Error(err))
		return nil, nil, err
	}

	logger.Info("DB connection successful")

	err = db.AutoMigrate(&domain.Users{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return nil, nil, err
	}

	cleanup = func() error {
		sqlDB, _ := db.DB()
		return sqlDB.Close()
	}

	return db, cleanup, nil
}
