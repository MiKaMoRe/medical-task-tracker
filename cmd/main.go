// Package main provides a CLI for application actions
//
// run this package only from the root project directory
//
// depends on .env or .env.local file
//
// read the README.md file or run:
//
//	go run ./cmd --help
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MiKaMoRe/medical-task-tracker/internal/config"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"
	"github.com/joho/godotenv"
)

func main() {
	log := logger.MustNewWithConfigLevel(config.EnvProd, config.LogLevelInfo)
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Printf("failed to close logger: %v\n", err)
		}
	}()

	// Check if no arguments
	if len(os.Args) < 2 {
		log.Warn("Command required")
		fmt.Println("No arguments")
		printHelp()
		os.Exit(1)
	}

	// Define CLI params
	command := os.Args[1]
	isLocal := flag.Bool("local", false, "Run application with .env.local environment")
	isHelp := flag.Bool("help", false, "Prints help")

	parseFlags(log)

	if *isHelp {
		printHelp()
		return
	}

	envFile := ".env"
	if *isLocal {
		envFile = ".env.local"
	}

	fmt.Printf("Running with %s\n", envFile)

	// Load environment variables
	err := godotenv.Load(envFile)
	if err != nil {
		log.Warn("Failed to load env file", "file", envFile, "error", err.Error())
		fmt.Printf("Error loading ENV file: %v", err)
	}

	// Call commands by name
	switch command {
	case "migration":
		migrateCommand()
	case "run":
		runCommand()
	default:
		log.Warn("Unknown command", "command", command)
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func parseFlags(log logger.Logger) {
	var flagsOnly []string
	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "-") {
			flagsOnly = append(flagsOnly, arg)
		}
	}

	if err := flag.CommandLine.Parse(flagsOnly); err != nil {
		log.Error("Failed to parse flags", "error", err.Error())
		fmt.Printf("Error parsing flags: %s\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd <command> [flags]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  run          Starts the main application application")
	fmt.Println("  migration    Runs database migrations (up, down, create)")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  --local      Run application with .env.local environment (default is .env)")
	fmt.Println("  --help   Show this help message")
	fmt.Println("")
	fmt.Println("Note: Run this package only from the root project directory.")
}
