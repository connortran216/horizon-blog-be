package main

import (
	"fmt"
	"go-crud/initializers"
	"go-crud/models"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	fmt.Println("Starting simple database migration...")

	// Auto-migrate all models
	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.PostVersion{},
		&models.Post{},
		&models.Tag{},
		&models.PostTag{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migration completed successfully!")
}
