package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
	"github.com/MiKaMoRe/medical-task-tracker/internal/db/migrate"
)

const migrationsDir = "internal/db/migrate/migrations"

func migrateCommand() {

	if len(os.Args) < 3 {
		fmt.Println("No arguments")
		printMigrationHelp()
		return
	}

	action := os.Args[2]

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	dsn := cfg.DB.PostgresDSN()

	switch action {
	case "up":
		fmt.Println("Running migrations...")
		if err := migrate.Up(dsn); err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully")
	case "down":
		fmt.Println("Rolling back migration...")
		if err := migrate.Down(dsn); err != nil {
			fmt.Printf("Rollback failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migration rolled back successfully")
	case "create":
		createMigration()
	default:
		fmt.Printf("Unknown action: %s\n", action)
		printMigrationHelp()
		os.Exit(1)
	}
}

func createMigration() {
	if len(os.Args) < 4 {
		fmt.Println("You need to specify migration name")
		printMigrationHelp()
		os.Exit(1)
	}

	if !commandExists("goose") {
		fmt.Println("Install goose CLI for local development:")
		fmt.Println("  go install github.com/pressly/goose/v3/cmd/goose@latest")
		os.Exit(1)
	}

	name := os.Args[3]
	cmd := exec.Command("goose", "-dir", migrationsDir, "create", name, "sql")
	output, err := cmd.CombinedOutput()
	fmt.Print(string(output))
	if err != nil {
		fmt.Printf("Command execution failed: %v\n", err)
		os.Exit(1)
	}
}

func printMigrationHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd migration <action> [name]")
	fmt.Println("")
	fmt.Println("Actions:")
	fmt.Println("  up       Apply pending migrations")
	fmt.Println("  down     Roll back the last migration")
	fmt.Println("  create   Create a new migration file (requires goose CLI, dev only)")
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
