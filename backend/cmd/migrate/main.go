package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dejavu/backend/pkg/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("✅ Migration completed successfully")
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		log.Println("✅ Rollback completed successfully")
	default:
		log.Fatal("Unknown command. Use 'up' or 'down'")
	}
}

func migrateUp(db *database.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name VARCHAR(255) NOT NULL,
			repo_url TEXT NOT NULL,
			build_command VARCHAR(255) DEFAULT 'npm run build',
			output_dir VARCHAR(255) DEFAULT 'dist',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS deployments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			status VARCHAR(50) DEFAULT 'pending',
			subdomain VARCHAR(255) UNIQUE NOT NULL,
			image_url TEXT,
			commit_hash VARCHAR(255),
			build_logs TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS billing_accounts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			credits DECIMAL(10, 2) DEFAULT 0.00,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS usage_records (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			deployment_id UUID REFERENCES deployments(id) ON DELETE SET NULL,
			type VARCHAR(50) NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_project_id ON deployments(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deployments_status ON deployments(status)`,
		`CREATE INDEX IF NOT EXISTS idx_usage_records_user_id ON usage_records(user_id)`,
	}

	for i, migration := range migrations {
		fmt.Printf("Running migration %d/%d...\n", i+1, len(migrations))
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

func migrateDown(db *database.DB) error {
	migrations := []string{
		`DROP TABLE IF EXISTS usage_records CASCADE`,
		`DROP TABLE IF EXISTS billing_accounts CASCADE`,
		`DROP TABLE IF EXISTS deployments CASCADE`,
		`DROP TABLE IF EXISTS projects CASCADE`,
		`DROP TABLE IF EXISTS users CASCADE`,
	}

	for i, migration := range migrations {
		fmt.Printf("Rolling back migration %d/%d...\n", i+1, len(migrations))
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("rollback %d failed: %w", i+1, err)
		}
	}

	return nil
}
