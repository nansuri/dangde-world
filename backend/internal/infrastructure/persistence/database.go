package persistence

import (
	"dangde-world/backend/internal/shared/config"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	if cfg.DatabaseURL != "" {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("dangde-world.db"), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&userModel{}, &categoryModel{}, &activityModel{}, &assignmentModel{}, &activityDataModel{}); err != nil {
		return nil, err
	}

	return db, nil
}
