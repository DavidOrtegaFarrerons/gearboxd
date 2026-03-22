package main

import (
	"bytes"
	"encoding/json"
	"gearboxd/internal/assert"
	"gearboxd/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
)

func TestCreateCarHandler(t *testing.T) {
	tests := []struct {
		Name        string          `json:"-"`
		ID          int64           `json:"id"`
		Make        string          `json:"make"`
		Model       string          `json:"model"`
		Year        int             `json:"year"`
		Description string          `json:"description"`
		ImageURL    string          `json:"image_url"`
		Gearbox     string          `json:"gearbox"`
		Drivetrain  string          `json:"drivetrain"`
		Horsepower  int             `json:"horsepower"`
		Fuel        string          `json:"fuel"`
		PriceNew    decimal.Decimal `json:"price_new"`
		Version     int             `json:"version"`
	}{
		{
			Name: "valid car",
			Make: "BMW", Model: "M3", Year: 2020,
			Description: "Sport sedan",
			ImageURL:    "https://img.com/m3.jpg",
			Gearbox:     "automatic", Drivetrain: "RWD",
			Horsepower: 480, Fuel: "gas", PriceNew: decimal.NewFromInt(50000),
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

			var responseCar data.Car
			err = json.NewDecoder(rr.Body).Decode(&responseCar)
			if err != nil {
				t.Errorf("got %v error", err)
			}

			assert.Equal(t, rr.Code, http.StatusCreated)
		})
	}
}
