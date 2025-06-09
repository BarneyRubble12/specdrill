package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	// Configure logrus
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetLevel(logrus.InfoLevel)
}

// TestCaseLog represents a structured log entry for a test case
type TestCaseLog struct {
	Name           string            `json:"name"`
	Endpoint       string            `json:"endpoint"`
	Method         string            `json:"method"`
	PathParams     map[string]string `json:"path_params"`
	QueryParams    map[string]string `json:"query_params"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestBody    string            `json:"request_body"`
	ResponseStatus int               `json:"response_status"`
	ResponseBody   string            `json:"response_body"`
	Error          string            `json:"error,omitempty"`
}

// LogTestCase logs a test case execution with all its details
func LogTestCase(testLog TestCaseLog) {
	fields := logrus.Fields{
		"name":            testLog.Name,
		"endpoint":        testLog.Endpoint,
		"method":          testLog.Method,
		"path_params":     testLog.PathParams,
		"query_params":    testLog.QueryParams,
		"request_headers": testLog.RequestHeaders,
		"request_body":    testLog.RequestBody,
	}

	if testLog.ResponseStatus != 0 {
		fields["response_status"] = testLog.ResponseStatus
		fields["response_body"] = testLog.ResponseBody
	}

	if testLog.Error != "" {
		log.WithFields(fields).Error("Test case failed")
	} else {
		log.WithFields(fields).Info("Test case executed")
	}
}
