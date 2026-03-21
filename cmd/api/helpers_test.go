package main

import (
	"gearboxd/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name         string
		appEnv       string
		envelope     envelope
		status       int
		header       http.Header
		expectedBody string
	}{
		{
			name:     "dev mode returns indented json",
			appEnv:   "dev",
			envelope: envelope{"status": "healthy"},
			status:   http.StatusOK,
			header:   nil,
			expectedBody: `{
	"status": "healthy"
}`,
		},
		{
			name:         "prod mode returns compact json",
			appEnv:       "prod",
			envelope:     envelope{"status": "healthy"},
			status:       http.StatusOK,
			header:       nil,
			expectedBody: `{"status":"healthy"}`,
		},
		{
			name:         "status is returned correctly",
			appEnv:       "prod",
			envelope:     envelope{"status": "healthy"},
			status:       http.StatusAccepted,
			header:       nil,
			expectedBody: `{"status":"healthy"}`,
		},
		{
			name:         "custom headers are returned correctly",
			appEnv:       "prod",
			envelope:     envelope{"status": "unauthorized"},
			status:       http.StatusUnauthorized,
			header:       http.Header{"X-Custom-Header": []string{"foo", "bar", "boz"}},
			expectedBody: `{"status":"unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			app := newTestApplication(t, &config{env: tt.appEnv})

			err := app.writeJSON(rr, tt.status, tt.envelope, tt.header)
			if err != nil {
				t.Errorf("Got error: %v", err)
			}

			assert.Equal(t, rr.Body.String(), tt.expectedBody)
			assert.Equal(t, rr.Code, tt.status)
			assert.Equal(t, rr.Header().Get("Content-Type"), "application/json")

			if tt.header != nil {
				for k, v := range tt.header {
					assert.Equal(t, rr.Header()[k], v)
				}
			}
		})
	}

}
