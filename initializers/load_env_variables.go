package initializers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
  	// Get the directory of this source file
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filename)) // Go up from initializers/ to project root
	envPath := filepath.Join(projectRoot, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}
}
