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

	"github.com/joho/godotenv"
)

func main() {
	// Check if no arguments
	if len(os.Args) < 2 {
		fmt.Println("No arguments")
		printHelp()
		os.Exit(1)
	}

	// Define CLI params
	command := os.Args[1]
	isLocal := flag.Bool("local", false, "Run application with .env.local environment")
	isHelp := flag.Bool("help", false, "Prints help")

	parseFlags()

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
		fmt.Printf("Error loading ENV file: %v", err)
	}

	// Call commands by name
	switch command {
	case "migration":
		migrateCommand()
	case "run":
		runCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func parseFlags() {
	var flagsOnly []string
	for _, arg := range os.Args[2:] {
		if strings.HasPrefix(arg, "-") {
			flagsOnly = append(flagsOnly, arg)
		}
	}

	if err := flag.CommandLine.Parse(flagsOnly); err != nil {
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
