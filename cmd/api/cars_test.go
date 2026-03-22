package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gearboxd/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCarHandler(t *testing.T) {
	tests := []struct {
		Name               string  `json:"-"`
		Make               string  `json:"make"`
		Model              string  `json:"model"`
		Year               int     `json:"year"`
		Description        string  `json:"description"`
		ImageURL           string  `json:"image_url"`
		Gearbox            string  `json:"gearbox"`
		Drivetrain         string  `json:"drivetrain"`
		Horsepower         int     `json:"horsepower"`
		Fuel               string  `json:"fuel"`
		PriceNew           float64 `json:"price_new"`
		ExpectedStatusCode int     `json:"-"`
		RandomField        string  `json:"random-field,omitempty"`
	}{
		{
			Name: "valid car",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			ExpectedStatusCode: 201,
		},
		{
			Name: "Extra fields not allowed",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			ExpectedStatusCode: 400,
			RandomField:        "I am a field",
		},
		{
			Name: "Invalid field returns validation error",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "does not exist", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: 50000,
			ExpectedStatusCode: 422,
		},
	}

	app := newTestApplication(t, nil)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			js, err := json.Marshal(tt)
			if err != nil {
				t.Errorf("got %v error", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/cars", bytes.NewReader(js))
			app.createCarHandler(rr, req)

			assert.Equal(t, rr.Code, tt.ExpectedStatusCode)
			if rr.Code == http.StatusCreated {
				assert.Equal(t, rr.Header().Get("Location"), fmt.Sprintf("/v1/cars/%d", 1))
			}
		})
	}
}
