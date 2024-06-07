package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-check-go/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

const MOZILLA_TLS_OBSERVATORY_API = "https://tls-observatory.services.mozilla.com/api/v1"

func TestTlsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		urlParam         string
		mockScanResp     string
		mockScanStatus   int
		mockResultResp   string
		mockResultStatus int
		expectedStatus   int
		expectedBody     map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "URL parameter is required"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "http://invalid-url",
			mockScanResp:   `{"scan_id": 0}`,
			mockScanStatus: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "Failed to get scan_id from TLS Observatory"},
		},
		{
			name:             "Valid URL with successful scan",
			urlParam:         "http://example.com",
			mockScanResp:     `{"scan_id": 12345}`,
			mockScanStatus:   http.StatusOK,
			mockResultResp:   `{"grade": "A+"}`,
			mockResultStatus: http.StatusOK,
			expectedStatus:   http.StatusOK,
			expectedBody:     map[string]interface{}{"grade": "A+"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()

			if tt.urlParam != "" {
				gock.New(MOZILLA_TLS_OBSERVATORY_API).
					Post("/scan").
					Reply(tt.mockScanStatus).
					BodyString(tt.mockScanResp)

				if tt.mockScanStatus == http.StatusOK && tt.mockResultResp != "" {
					gock.New(MOZILLA_TLS_OBSERVATORY_API).
						Get("/results").
						MatchParam("id", "12345").
						Reply(tt.mockResultStatus).
						BodyString(tt.mockResultResp)
				}
			}

			router := gin.Default()
			router.GET("/tls", func(c *gin.Context) {
				ctrl := &controllers.TlsController{}
				ctrl.TlsHandler(c)
			})

			req, _ := http.NewRequest("GET", "/tls?url="+tt.urlParam, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestHandleTLS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		urlParam         string
		mockScanResp     string
		mockScanStatus   int
		mockResultResp   string
		mockResultStatus int
		expectedStatus   int
		expectedBody     map[string]interface{}
	}{
		{
			name:           "Missing URL parameter",
			urlParam:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "missing URL parameter"},
		},
		{
			name:           "Invalid URL",
			urlParam:       "http://invalid-url",
			mockScanResp:   `{"scan_id": 0}`,
			mockScanStatus: http.StatusOK,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to get scan_id from TLS Observatory"},
		},
		{
			name:             "Valid URL with successful scan",
			urlParam:         "http://example.com",
			mockScanResp:     `{"scan_id": 12345}`,
			mockScanStatus:   http.StatusOK,
			mockResultResp:   `{"grade": "A+"}`,
			mockResultStatus: http.StatusOK,
			expectedStatus:   http.StatusOK,
			expectedBody:     map[string]interface{}{"grade": "A+"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			defer gock.Off()

			if tc.urlParam != "" {
				gock.New(MOZILLA_TLS_OBSERVATORY_API).
					Post("/scan").
					Reply(tc.mockScanStatus).
					BodyString(tc.mockScanResp)

				if tc.mockScanStatus == http.StatusOK && tc.mockResultResp != "" {
					gock.New(MOZILLA_TLS_OBSERVATORY_API).
						Get("/results").
						MatchParam("id", "12345").
						Reply(tc.mockResultStatus).
						BodyString(tc.mockResultResp)
				}
			}

			req := httptest.NewRequest("GET", "/tls?url="+tc.urlParam, nil)
			rec := httptest.NewRecorder()
			controllers.HandleTLS().ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var responseBody map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)
		})
	}
}
