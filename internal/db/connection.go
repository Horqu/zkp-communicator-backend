package db

import (
	"fmt"
	// "os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectGORM() (*gorm.DB, error) {
	// dsn := fmt.Sprintf(
	// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASSWORD"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("DB_PORT"),
	// )

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost",
		"zkp_user",
		"zkp_password",
		"zkp_db",
		"5432",
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Session{},
		&Contact{},
		&Message{},
		&Group{},
		&GroupMember{},
		&AgeProof{},
	)
}
