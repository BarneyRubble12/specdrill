package parser

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSpec(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "specdrill-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test spec file
	specContent := `{
		"openapi": "3.0.0",
		"info": {
			"title": "Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/test": {
				"get": {
					"responses": {
						"200": {
							"description": "OK"
						}
					}
				}
			}
		}
	}`
	specPath := filepath.Join(tempDir, "test-spec.json")
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
		t.Fatalf("Failed to write test spec: %v", err)
	}

	// Create a test server for remote spec
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(specContent))
	}))
	defer server.Close()

	tests := []struct {
		name     string
		specPath string
		baseURL  string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Local spec with base URL",
			specPath: specPath,
			baseURL:  "http://localhost:8080",
			wantErr:  false,
		},
		{
			name:     "Local spec without base URL",
			specPath: specPath,
			baseURL:  "",
			wantErr:  true,
			errMsg:   "base URL is required for local spec files",
		},
		{
			name:     "Remote spec with base URL",
			specPath: server.URL,
			baseURL:  "http://api.example.com",
			wantErr:  false,
		},
		{
			name:     "Remote spec without base URL",
			specPath: server.URL,
			baseURL:  "",
			wantErr:  false,
		},
		{
			name:     "Invalid local spec path",
			specPath: "nonexistent.json",
			baseURL:  "http://localhost:8080",
			wantErr:  true,
			errMsg:   "failed to read spec file",
		},
		{
			name:     "Invalid remote spec URL",
			specPath: "http://invalid-url",
			baseURL:  "http://localhost:8080",
			wantErr:  true,
			errMsg:   "failed to fetch spec from URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			spec, err := parser.ParseSpec(tt.specPath, tt.baseURL)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, spec)
			assert.Equal(t, "3.0.0", spec.OpenAPI)
			assert.Equal(t, "Test API", spec.Info.Title)
			assert.Equal(t, "1.0.0", spec.Info.Version)

			// Check base URL
			if tt.baseURL != "" {
				assert.Equal(t, tt.baseURL, spec.BaseURL)
			} else if strings.HasPrefix(tt.specPath, "http") {
				// For remote specs without base URL, it should be derived from the spec URL
				assert.Equal(t, server.URL, spec.BaseURL)
			}

			// Check paths
			assert.NotNil(t, spec.Paths)
			assert.Contains(t, spec.Paths, "/test")
			assert.NotNil(t, spec.Paths["/test"].Get)
			assert.Contains(t, spec.Paths["/test"].Get.Responses, "200")
		})
	}
}
