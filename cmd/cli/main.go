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
	spec, err := container.Parser.ParseSpec(*specPath, *baseURL)
	if err != nil {
		fmt.Printf("Error parsing spec: %v\n", err)
		os.Exit(1)
	}

	totalTests := 0
	passedTests := 0
	failedTests := 0

	fmt.Printf("\nTest Results for API (Base URL: %s)\n", spec.BaseURL)

	for path, pathItem := range spec.Paths {
		endpoints := []struct {
			method string
			op     interface{}
		}{
			{"GET", pathItem.Get},
			{"POST", pathItem.Post},
			{"PUT", pathItem.Put},
			{"DELETE", pathItem.Delete},
			{"PATCH", pathItem.Patch},
			{"OPTIONS", pathItem.Options},
			{"HEAD", pathItem.Head},
		}
		for _, ep := range endpoints {
			if ep.op == nil {
				continue
			}
			totalTests++
			result, err := container.Executor.ExecuteTest(spec, path, ep.method)
			if err != nil {
				failedTests++
				fmt.Printf("✗ %s %s\n  Error: %v\n", ep.method, path, err)
				continue
			}
			if result.StatusCode >= 200 && result.StatusCode < 300 {
				passedTests++
				fmt.Printf("✓ %s %s (%d)\n", ep.method, path, result.StatusCode)
			} else {
				failedTests++
				fmt.Printf("✗ %s %s (%d)\n  Response: %s\n", ep.method, path, result.StatusCode, result.Body)
			}
		}
	}

	fmt.Printf("\nTotal Tests: %d\n", totalTests)
	fmt.Printf("Passed: %d\n", passedTests)
	fmt.Printf("Failed: %d\n", failedTests)

	// Exit with non-zero status if any tests failed
	if failedTests > 0 {
		os.Exit(1)
	}
}
