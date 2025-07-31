package main

import (
	"fmt"
	"os"

	"github.com/gerrytan/azsubsyn/internal/apply"
	"github.com/gerrytan/azsubsyn/internal/credential"
	"github.com/gerrytan/azsubsyn/internal/plan"
)

var Version = "dev-build"
var GitCommitSHA = "unknown"
var BuildNumber = "unknown"
var BuildDate = "unknown"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "credcheck":
		if err := credential.RunCredCheck(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "plan":
		if err := plan.RunPlan(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "apply":
		if err := apply.RunApply(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "version", "-v", "--version":
		fmt.Printf("version: %s\ngit commit SHA: %s\nbuild number: %s\nbuild date: %s\n",
			Version, GitCommitSHA, BuildNumber, BuildDate)
		os.Exit(0)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("azsubsyn - A CLI tool to ensure target Azure subscription has all RPs (resource providers) and ")
	fmt.Println("           preview features registered compared to source (which can be on a different tenant).")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  azsubsyn <command>")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  credcheck    Check credentials and connectivity to both source and target subscriptions")
	fmt.Println("  plan         Scan unregistered RPs and preview feature in the target subscription and save the plan to a file")
	fmt.Println("  apply        Apply the plan file to the target subscription")
	fmt.Println("  help         Show this help message")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  See https://github.com/gerrytan/azsubsyn for credential setup and usage example.")
}
