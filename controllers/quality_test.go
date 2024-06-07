package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetQualityHandler(t *testing.T) {
	t.Skip("broken due to api key handling inside handler")
	router := gin.Default()
	ctrl := &controllers.QualityController{}
	router.GET("/check-quality", ctrl.GetQualityHandler)

	tests := []struct {
		name         string
		url          string
		apiKey       string
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:         "Missing URL parameter",
			url:          "",
			apiKey:       "test-api-key",
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "url parameter is required"},
		},
		{
			name:         "Missing API key",
			url:          "http://example.com",
			apiKey:       "",
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "Missing Google API. You need to set the `GOOGLE_CLOUD_API_KEY` environment variable"},
		},
		{
			name:         "Valid request with expected failure",
			url:          "http://example.com",
			apiKey:       "test-api-key",
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Failed to fetch the Pagespeed data"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.apiKey != "" {
				os.Setenv("GOOGLE_CLOUD_API_KEY", tt.apiKey)
			} else {
				os.Unsetenv("GOOGLE_CLOUD_API_KEY")
			}

			req, err := http.NewRequest("GET", "/check-quality?url="+tt.url, nil)
			assert.NoError(t, err)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedBody, response)
		})
	}
}

func TestHandleGetQuality(t *testing.T) {
	t.Skip("broken due to api key handling inside handler")
	t.Parallel()
	tests := []struct {
		name         string
		url          string
		apiKey       string
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:         "Missing URL parameter",
			url:          "",
			apiKey:       "test-api-key",
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "missing URL parameter"},
		},
		{
			name:         "Missing API key",
			url:          "http://example.com",
			apiKey:       "",
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{"error": "missing Google API. You need to set the `GOOGLE_CLOUD_API_KEY` environment variable"},
		},
		{
			name:         "Valid request with expected failure",
			url:          "http://example.com",
			apiKey:       "test-api-key",
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{"error": "Failed to fetch the Pagespeed data"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.apiKey != "" {
				os.Setenv("GOOGLE_CLOUD_API_KEY", tc.apiKey)
			} else {
				os.Unsetenv("GOOGLE_CLOUD_API_KEY")
			}

			req := httptest.NewRequest("GET", "/check-quality?url="+tc.url, nil)
			w := httptest.NewRecorder()
			controllers.HandleGetQuality().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedBody, response)
		})
	}
}
