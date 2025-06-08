package parser

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/BarneyRubble12/specdrill/internal/core/model"
	"github.com/getkin/kin-openapi/openapi3"
)

// Parser defines the interface for parsing OpenAPI specifications
type Parser interface {
	ParseSpec(specPath string) (*model.TestSuite, error)
}

// OpenAPIParser implements the Parser interface for OpenAPI specifications
type OpenAPIParser struct {
	client *http.Client
}

// NewOpenAPIParser creates a new OpenAPI parser
func NewOpenAPIParser() *OpenAPIParser {
	return &OpenAPIParser{
		client: &http.Client{},
	}
}

// ParseSpec parses an OpenAPI specification file or URL and returns a TestSuite
func (p *OpenAPIParser) ParseSpec(specPath string) (*model.TestSuite, error) {
	// Read the spec data
	data, err := p.readSpecData(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec: %w", err)
	}

	// Load the OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Create a new test suite
	suite := &model.TestSuite{
		Name:    doc.Info.Title,
		BaseURL: doc.Servers[0].URL, // Use the first server URL as base URL
	}

	// Generate test cases for each path and method
	for path, pathItem := range doc.Paths.Map() {
		if pathItem.Get != nil {
			testCase := createTestCase("GET", path, pathItem.Get)
			suite.TestCases = append(suite.TestCases, testCase)
		}
		if pathItem.Post != nil {
			testCase := createTestCase("POST", path, pathItem.Post)
			suite.TestCases = append(suite.TestCases, testCase)
		}
		if pathItem.Put != nil {
			testCase := createTestCase("PUT", path, pathItem.Put)
			suite.TestCases = append(suite.TestCases, testCase)
		}
		if pathItem.Delete != nil {
			testCase := createTestCase("DELETE", path, pathItem.Delete)
			suite.TestCases = append(suite.TestCases, testCase)
		}
	}

	return suite, nil
}

// readSpecData reads the OpenAPI specification from either a local file or a URL
func (p *OpenAPIParser) readSpecData(specPath string) ([]byte, error) {
	if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
		// Fetch from URL
		resp, err := p.client.Get(specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch spec from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch spec: HTTP %d", resp.StatusCode)
		}

		return io.ReadAll(resp.Body)
	}

	// Read from local file
	return os.ReadFile(specPath)
}

// createTestCase creates a test case from an operation
func createTestCase(method, path string, operation *openapi3.Operation) model.TestCase {
	testCase := model.TestCase{
		Name:           fmt.Sprintf("%s %s", method, path),
		Method:         method,
		Path:           path,
		ExpectedStatus: getExpectedStatus(operation),
		Description:    operation.Summary,
	}

	// Add request body if present
	if operation.RequestBody != nil {
		// TODO: Extract example request body from schema
		testCase.RequestBody = nil
	}

	return testCase
}

// getExpectedStatus extracts the expected status code from the operation
func getExpectedStatus(operation *openapi3.Operation) int {
	// Default to 200 if no responses are defined
	if operation.Responses == nil {
		return 200
	}

	// Look for 2xx responses first
	for status := range operation.Responses.Map() {
		if status[0] == '2' {
			return 200 // Default to 200 for any 2xx response
		}
	}

	// Just return 200 as default since we can't determine the actual status
	return 200
}
