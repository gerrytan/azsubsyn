package main

import (
	"fmt"
	"os"

	"github.com/gerrytan/azdiffit/internal/apply"
	"github.com/gerrytan/azdiffit/internal/credential"
	"github.com/gerrytan/azdiffit/internal/plan"
)

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
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("azdiffit - A CLI tool to setup a target Azure subscription based on a source subscription")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  azdiffit <command>")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  credcheck    Check credentials and connectivity to both source and target subscriptions")
	fmt.Println("  help         Show this help message")
	fmt.Println()
	fmt.Println("DESCRIPTION:")
	fmt.Println("  azdiffit helps setup a target Azure subscription based on a source subscription.")
	fmt.Println("  It ensures the target subscription has all Resource Providers registered,")
	fmt.Println("  preview features registered, and enough quotas.")
}
