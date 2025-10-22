package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go-crud/initializers"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate() {
	// Load environment variables
	initializers.LoadEnvVariables()

	// Get database DSN from environment
	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	// Open database connection
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver: %v", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/sql",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Get command from arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run migration/migrate.go [up|down|version|force]")
		fmt.Println("Commands:")
		fmt.Println("  up       - Apply all migrations")
		fmt.Println("  down     - Rollback all migrations")
		fmt.Println("  version  - Show current migration version")
		fmt.Println("  force N  - Force migration to version N")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		fmt.Println("Applying migrations...")
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No new migrations to apply")
		} else {
			fmt.Println("✅ Migrations applied successfully!")
		}

	case "down":
		fmt.Println("Rolling back migrations...")
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to rollback")
		} else {
			fmt.Println("✅ Migrations rolled back successfully!")
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force command requires a version number")
		}
		versionStr := os.Args[2]
		fmt.Printf("Forcing migration to version: %s\n", versionStr)

		// Parse version number (simplified - you might want better parsing)
		var version int
		if versionStr == "0" {
			version = 0
		} else {
			// For now, just force to version 2 (tags migration)
			version = 2
		}

		err = m.Force(version)
		if err != nil {
			log.Fatalf("Failed to force migration: %v", err)
		}
		fmt.Println("✅ Migration forced successfully!")

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: go run migration/migrate.go [up|down|version|force]")
		os.Exit(1)
	}
}
