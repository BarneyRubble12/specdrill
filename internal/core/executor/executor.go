package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/BarneyRubble12/specdrill/internal/core/domain"
	"github.com/BarneyRubble12/specdrill/internal/core/logger"
)

// Executor handles the execution of API tests
type Executor struct {
	client *http.Client
}

// NewExecutor creates a new Executor instance
func NewExecutor() *Executor {
	return &Executor{
		client: &http.Client{},
	}
}

// ExecuteTest runs a test case against the API
func (e *Executor) ExecuteTest(spec *domain.APISpec, path string, method string) (*TestResult, error) {
	// Construct the full URL
	baseURL := spec.BaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Replace path parameters with actual values
	path = replacePathParams(path)

	fullURL := baseURL + path

	// Create the request
	var req *http.Request
	var err error

	// Analyze the spec to determine if a request body is required
	// For simplicity, we'll assume POST/PUT/PATCH methods require a body
	var body []byte
	if method == "POST" || method == "PUT" || method == "PATCH" {
		// Generate a simple JSON body for demonstration
		body = []byte(`{"name": "test"}`)
	}

	if body != nil {
		req, err = http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, fullURL, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Create test case log
	testLog := logger.TestCaseLog{
		Name:           fmt.Sprintf("%s %s", method, path),
		Endpoint:       fullURL,
		Method:         method,
		PathParams:     extractPathParams(path),
		QueryParams:    extractQueryParams(req.URL),
		RequestHeaders: extractHeaders(req.Header),
		RequestBody:    string(body),
	}

	// Execute the request
	resp, err := e.client.Do(req)
	if err != nil {
		testLog.Error = err.Error()
		logger.LogTestCase(testLog)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		testLog.Error = err.Error()
		logger.LogTestCase(testLog)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the test result
	result := &TestResult{
		URL:        fullURL,
		Method:     method,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       string(respBody),
	}

	// Check if the response is valid JSON
	var jsonBody interface{}
	if err := json.Unmarshal(respBody, &jsonBody); err != nil {
		result.IsValidJSON = false
	} else {
		result.IsValidJSON = true
	}

	// Update and log test case result
	testLog.ResponseStatus = resp.StatusCode
	testLog.ResponseBody = string(respBody)
	logger.LogTestCase(testLog)

	return result, nil
}

// replacePathParams replaces path parameters with actual values
func replacePathParams(path string) string {
	// For demonstration, replace {id} with a sample value
	return strings.Replace(path, "{id}", "1", -1)
}

// extractPathParams extracts path parameters from the path
func extractPathParams(path string) map[string]string {
	params := make(map[string]string)
	if strings.Contains(path, "{id}") {
		params["id"] = "1"
	}
	return params
}

// extractQueryParams extracts query parameters from the URL
func extractQueryParams(u *url.URL) map[string]string {
	params := make(map[string]string)
	q := u.Query()
	for k := range q {
		params[k] = q.Get(k)
	}
	return params
}

// extractHeaders extracts headers from the request
func extractHeaders(h http.Header) map[string]string {
	headers := make(map[string]string)
	for k, v := range h {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	return headers
}

// TestResult represents the result of a test execution
type TestResult struct {
	URL         string
	Method      string
	StatusCode  int
	Headers     http.Header
	Body        string
	IsValidJSON bool
}
