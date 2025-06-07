package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hrpd/specdrill/internal/core/executor"
	"github.com/hrpd/specdrill/internal/core/parser"
)

func main() {
	// Parse command line flags
	specPath := flag.String("spec", "", "Path to OpenAPI specification file (YAML or JSON)")
	baseURL := flag.String("base-url", "", "Base URL for the API (overrides server URL in spec)")
	flag.Parse()

	// Validate required flags
	if *specPath == "" {
		fmt.Println("Error: --spec flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create parser and executor
	parser := parser.NewOpenAPIParser()
	executor := executor.NewHTTPExecutor()

	// Parse the OpenAPI spec
	suite, err := parser.ParseSpec(*specPath)
	if err != nil {
		fmt.Printf("Error parsing spec: %v\n", err)
		os.Exit(1)
	}

	// Override base URL if provided
	if *baseURL != "" {
		suite.BaseURL = *baseURL
	}

	// Execute the test suite
	summary := executor.ExecuteSuite(suite)

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