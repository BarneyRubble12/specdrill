package executor

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BarneyRubble12/specdrill/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestExecuteTest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request headers
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Return a test response
		response := map[string]interface{}{
			"message": "test response",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create a test spec
	spec := &domain.APISpec{
		BaseURL: server.URL,
	}

	tests := []struct {
		name     string
		path     string
		method   string
		wantErr  bool
		errMsg   string
		validate func(*testing.T, *TestResult)
	}{
		{
			name:    "Successful GET request",
			path:    "/test",
			method:  "GET",
			wantErr: false,
			validate: func(t *testing.T, result *TestResult) {
				assert.Equal(t, http.StatusOK, result.StatusCode)
				assert.True(t, result.IsValidJSON)
				assert.Contains(t, result.Body, "test response")
			},
		},
		{
			name:    "Request with leading slash",
			path:    "/test",
			method:  "GET",
			wantErr: false,
			validate: func(t *testing.T, result *TestResult) {
				assert.Equal(t, http.StatusOK, result.StatusCode)
				assert.True(t, result.IsValidJSON)
			},
		},
		{
			name:    "Request without leading slash",
			path:    "test",
			method:  "GET",
			wantErr: false,
			validate: func(t *testing.T, result *TestResult) {
				assert.Equal(t, http.StatusOK, result.StatusCode)
				assert.True(t, result.IsValidJSON)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor()
			result, err := executor.ExecuteTest(spec, tt.path, tt.method)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.method, result.Method)
			assert.Equal(t, "application/json", result.Headers.Get("Content-Type"))

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
