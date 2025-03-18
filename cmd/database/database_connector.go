package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/Horqu/zkp-communicator-backend/cmd/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func GetDBInstance() *gorm.DB {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			"localhost",
			"zkp_user",
			"zkp_password",
			"zkp_db",
			"5432",
		)

		var err error
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		err = AutoMigrateAll(dbInstance)
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}

		log.Println("Database connection established and migrations applied.")
	})

	return dbInstance
}

func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Contact{},
		&models.Message{},
	)
}
