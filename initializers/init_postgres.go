package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	connection := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database: ", err)
		return
	}
	DB = db
}
