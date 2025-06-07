package model

import (
	"net/http"
)

// TestCase represents a single API test case
type TestCase struct {
	Name           string
	Method         string
	Path           string
	RequestBody    interface{}
	ExpectedStatus int
	Description    string
}

// TestResult represents the result of executing a test case
type TestResult struct {
	TestCase    TestCase
	Success     bool
	StatusCode  int
	Error       error
	Response    *http.Response
	Duration    int64 // in milliseconds
	Message     string
}

// TestSuite represents a collection of test cases
type TestSuite struct {
	Name      string
	BaseURL   string
	TestCases []TestCase
}

// TestSummary represents the summary of test execution
type TestSummary struct {
	TotalTests  int
	PassedTests int
	FailedTests int
	Duration    int64 // in milliseconds
	Results     []TestResult
} 