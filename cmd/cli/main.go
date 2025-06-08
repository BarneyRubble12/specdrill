package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BarneyRubble12/specdrill/internal/di"
)

func main() {
	// Parse command line flags
	specPath := flag.String("spec", "", "Path to OpenAPI specification file (YAML/JSON) or URL")
	baseURL := flag.String("base-url", "", "Base URL for the API (overrides server URL in spec)")
	flag.Parse()

	// Validate required flags
	if *specPath == "" {
		fmt.Println("Error: --spec flag is required")
		fmt.Println("Usage: specdrill --spec <file-path-or-url> [--base-url <api-base-url>]")
		fmt.Println("\nExamples:")
		fmt.Println("  specdrill --spec ./openapi.yaml")
		fmt.Println("  specdrill --spec https://api.example.com/openapi.json")
		fmt.Println("  specdrill --spec ./openapi.yaml --base-url https://staging-api.example.com")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize the application container
	container, err := di.InitializeContainer()
	if err != nil {
		fmt.Printf("Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Parse the OpenAPI spec
	suite, err := container.Parser.ParseSpec(*specPath)
	if err != nil {
		fmt.Printf("Error parsing spec: %v\n", err)
		os.Exit(1)
	}

	// Override base URL if provided
	if *baseURL != "" {
		suite.BaseURL = *baseURL
	}

	// Execute the test suite
	summary := container.Executor.ExecuteSuite(suite)

	// Print results
	fmt.Printf("\nTest Results for %s\n", suite.Name)
	fmt.Printf("Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("Passed: %d\n", summary.PassedTests)
	fmt.Printf("Failed: %d\n", summary.FailedTests)
	fmt.Printf("Duration: %dms\n\n", summary.Duration)

	// Print individual test results
	for _, result := range summary.Results {
		status := "✓"
		if !result.Success {
			status = "✗"
		}
		fmt.Printf("%s %s (%dms)\n", status, result.TestCase.Name, result.Duration)
		if !result.Success {
			fmt.Printf("  Error: %s\n", result.Message)
		}
	}

	// Exit with non-zero status if any tests failed
	if summary.FailedTests > 0 {
		os.Exit(1)
	}
}
