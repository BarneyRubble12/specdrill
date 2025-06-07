package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hrpd/specdrill/internal/core/model"
)

// Executor defines the interface for executing test cases
type Executor interface {
	ExecuteTest(testCase model.TestCase, baseURL string) model.TestResult
	ExecuteSuite(suite *model.TestSuite) model.TestSummary
}

// HTTPExecutor implements the Executor interface using HTTP requests
type HTTPExecutor struct {
	client *http.Client
}

// NewHTTPExecutor creates a new HTTP executor
func NewHTTPExecutor() *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ExecuteTest executes a single test case and returns the result
func (e *HTTPExecutor) ExecuteTest(testCase model.TestCase, baseURL string) model.TestResult {
	startTime := time.Now()
	result := model.TestResult{
		TestCase: testCase,
	}

	// Construct the full URL
	url := fmt.Sprintf("%s%s", baseURL, testCase.Path)

	// Create the request
	var req *http.Request
	var err error

	if testCase.RequestBody != nil {
		body, err := json.Marshal(testCase.RequestBody)
		if err != nil {
			result.Error = fmt.Errorf("failed to marshal request body: %w", err)
			result.Success = false
			return result
		}
		req, err = http.NewRequest(testCase.Method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(testCase.Method, url, nil)
	}

	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %w", err)
		result.Success = false
		return result
	}

	// Execute the request
	resp, err := e.client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("request failed: %w", err)
		result.Success = false
		return result
	}
	defer resp.Body.Close()

	// Record the result
	result.StatusCode = resp.StatusCode
	result.Response = resp
	result.Duration = time.Since(startTime).Milliseconds()
	result.Success = resp.StatusCode == testCase.ExpectedStatus

	if !result.Success {
		result.Message = fmt.Sprintf("Expected status %d but got %d", testCase.ExpectedStatus, resp.StatusCode)
	}

	return result
}

// ExecuteSuite executes all test cases in a suite and returns a summary
func (e *HTTPExecutor) ExecuteSuite(suite *model.TestSuite) model.TestSummary {
	summary := model.TestSummary{
		TotalTests: len(suite.TestCases),
	}

	startTime := time.Now()

	for _, testCase := range suite.TestCases {
		result := e.ExecuteTest(testCase, suite.BaseURL)
		summary.Results = append(summary.Results, result)

		if result.Success {
			summary.PassedTests++
		} else {
			summary.FailedTests++
		}
	}

	summary.Duration = time.Since(startTime).Milliseconds()
	return summary
} 