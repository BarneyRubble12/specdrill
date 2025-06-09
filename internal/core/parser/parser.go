package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/BarneyRubble12/specdrill/internal/core/domain"
	"github.com/BarneyRubble12/specdrill/internal/core/model"
	"github.com/getkin/kin-openapi/openapi3"
)

// Parser handles parsing of OpenAPI specifications
type Parser struct {
	client *http.Client
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{
		client: &http.Client{},
	}
}

// ParseSpec parses an OpenAPI specification from a file or URL
func (p *Parser) ParseSpec(specPath string, baseURL string) (*domain.APISpec, error) {
	var specData []byte
	var err error

	// Check if specPath is a URL
	if strings.HasPrefix(specPath, "http://") || strings.HasPrefix(specPath, "https://") {
		// For remote specs, if no base URL is provided, derive it from the spec URL
		if baseURL == "" {
			parsedURL, err := url.Parse(specPath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse spec URL: %w", err)
			}
			// Remove the path and query components to get the base URL
			parsedURL.Path = ""
			parsedURL.RawQuery = ""
			baseURL = parsedURL.String()
		}

		// Fetch the spec from URL
		resp, err := p.client.Get(specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch spec from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to fetch spec: HTTP %d", resp.StatusCode)
		}

		specData, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read spec from URL: %w", err)
		}
	} else {
		// For local files, if no base URL is provided, return an error
		if baseURL == "" {
			return nil, fmt.Errorf("base URL is required for local spec files")
		}

		// Read the spec from file
		specData, err = os.ReadFile(specPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read spec file: %w", err)
		}
	}

	// Parse the spec
	var spec domain.APISpec
	if err := json.Unmarshal(specData, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}

	// Set the base URL
	spec.BaseURL = baseURL

	return &spec, nil
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
