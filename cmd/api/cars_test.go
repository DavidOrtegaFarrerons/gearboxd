package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gearboxd/internal/assert"
	"gearboxd/internal/data"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
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

	app := newTestApplication(t, nil, nil)

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

func TestGetCarHandler(t *testing.T) {
	tests := []struct {
		name         string
		ID           int
		expectedCode int
	}{
		{
			name:         "Returns a car",
			ID:           1,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Returns 404 not found",
			ID:           2,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid ID",
			ID:           0,
			expectedCode: http.StatusNotFound,
		},
	}

	cars := []data.Car{
		{
			ID:          1,
			Make:        "BMW",
			Model:       "M3 Competition",
			Year:        2022,
			Description: "High-performance sports sedan with twin-turbo inline-6 engine",
			ImageURL:    "https://images.unsplash.com/photo-1619767886558-efdc259cde1a",
			Gearbox:     "automatic",
			Drivetrain:  "RWD",
			Horsepower:  510,
			Fuel:        "gas",
			PriceNew:    decimal.NewFromInt(85000),
			Version:     1,
		},
	}

	models := &data.Models{
		Cars: &MockCarModel{cars: cars},
	}

	app := newTestApplication(t, nil, models)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := createTestRequestWithIdParam(t, http.MethodGet, "/v1/cars", tt.ID, nil)
			app.getCarHandler(rr, req)

			assert.Equal(t, rr.Code, tt.expectedCode)
		})
	}
}

func TestDeleteCarHandler(t *testing.T) {
	tests := []struct {
		name         string
		ID           int
		expectedCode int
	}{
		{
			name:         "Deletes a car",
			ID:           1,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Returns 404 not found",
			ID:           7,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Invalid ID",
			ID:           0,
			expectedCode: http.StatusNotFound,
		},
	}

	cars := []data.Car{
		{
			ID: 1,
		},
		{
			ID: 2,
		},
	}

	models := &data.Models{
		Cars: &MockCarModel{cars: cars},
	}

	app := newTestApplication(t, nil, models)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			req := createTestRequestWithIdParam(t, http.MethodGet, "/v1/cars", tt.ID, nil)

			app.deleteCarHandler(rr, req)

			assert.Equal(t, rr.Code, tt.expectedCode)
		})
	}
}

func TestUpdateCarHandler(t *testing.T) {
	tests := []struct {
		name          string
		ID            int
		body          map[string]any
		versionHeader string
		expectedCode  int
	}{

		{
			name: "Updates make successfully",
			ID:   1,
			body: map[string]any{
				"make": "Audi",
			},
			versionHeader: "1",
			expectedCode:  http.StatusOK,
		},
		{
			name: "Partial update multiple fields",
			ID:   1,
			body: map[string]any{
				"make":  "Audi",
				"model": "RS5",
			},
			versionHeader: "1",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "No fields provided",
			ID:            1,
			body:          map[string]any{},
			versionHeader: "1",
			expectedCode:  http.StatusOK,
		},
		{
			name: "Invalid gearbox value",
			ID:   1,
			body: map[string]any{
				"gearbox": "invalid",
			},
			versionHeader: "1",
			expectedCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid fuel type",
			ID:   1,
			body: map[string]any{
				"fuel": "water",
			},
			versionHeader: "1",
			expectedCode:  http.StatusUnprocessableEntity,
		},
		{
			name: "Negative horsepower",
			ID:   1,
			body: map[string]any{
				"horsepower": -100,
			},
			versionHeader: "1",
			expectedCode:  http.StatusUnprocessableEntity,
		},
		{
			name:          "Malformed JSON",
			ID:            1,
			body:          nil,
			versionHeader: "1",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name: "Missing version header (should still update)",
			ID:   1,
			body: map[string]any{
				"make": "Audi",
			},
			versionHeader: "",
			expectedCode:  http.StatusOK,
		},
		{
			name: "Conflict due to wrong version",
			ID:   1,
			body: map[string]any{
				"make": "Audi",
			},
			versionHeader: "999",
			expectedCode:  http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cars := []data.Car{
				{
					ID:          1,
					Make:        "BMW",
					Model:       "M3 Competition",
					Year:        2022,
					Description: "High-performance sports sedan with twin-turbo inline-6 engine",
					ImageURL:    "https://images.unsplash.com/photo-1619767886558-efdc259cde1a",
					Gearbox:     "automatic",
					Drivetrain:  "RWD",
					Horsepower:  510,
					Fuel:        "gas",
					PriceNew:    decimal.NewFromInt(85000),
					Version:     1,
				},
				{
					ID:          2,
					Make:        "Toyota",
					Model:       "Corolla Hybrid",
					Year:        2023,
					Description: "Efficient and reliable hybrid compact sedan for daily driving",
					ImageURL:    "https://images.unsplash.com/photo-1606664515524-ed2f786a0bd6",
					Gearbox:     "automatic",
					Drivetrain:  "FWD",
					Horsepower:  140,
					Fuel:        "HEV",
					PriceNew:    decimal.NewFromInt(28000),
					Version:     1,
				},
			}

			models := &data.Models{
				Cars: &MockCarModel{cars: cars},
			}

			app := newTestApplication(t, nil, models)

			rr := httptest.NewRecorder()

			var bodyReader io.Reader

			if tt.body != nil {
				js, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatal(err)
				}
				bodyReader = bytes.NewReader(js)
			} else {
				bodyReader = bytes.NewReader([]byte("{invalid-json"))
			}

			req := createTestRequestWithIdParam(t, http.MethodPatch, "/v1/cars", tt.ID, bodyReader)
			req.Header.Set("Content-Type", "application/json")

			if tt.versionHeader != "" {
				req.Header.Set("X-Expected-Version", tt.versionHeader)
			}

			app.updateCarHandler(rr, req)

			assert.Equal(t, rr.Code, tt.expectedCode)
		})
	}
}

func TestListCarsHandler(t *testing.T) {
	cars := []data.Car{
		{ID: 1, Make: "BMW", Model: "M3", Year: 2020, Gearbox: "automatic", Drivetrain: "RWD", Horsepower: 480, Fuel: "gas", PriceNew: decimal.NewFromInt(50000)},
		{ID: 2, Make: "Audi", Model: "RS5", Year: 2021, Gearbox: "automatic", Drivetrain: "AWD", Horsepower: 450, Fuel: "gas", PriceNew: decimal.NewFromInt(60000)},
		{ID: 3, Make: "Volvo", Model: "XC60", Year: 2023, Gearbox: "automatic", Drivetrain: "AWD", Horsepower: 250, Fuel: "hybrid", PriceNew: decimal.NewFromInt(52000)},
	}

	models := &data.Models{
		Cars: &MockCarModel{cars: cars},
	}

	app := newTestApplication(t, nil, models)

	tests := []struct {
		name         string
		query        string
		expectedCode int
		expectedLen  int
		assertFn     func(t *testing.T, cars []data.Car)
	}{
		{
			name:         "Returns all cars",
			query:        "/v1/cars",
			expectedCode: http.StatusOK,
			expectedLen:  3,
		},
		{
			name:         "Filter by make",
			query:        "/v1/cars?make=BMW",
			expectedCode: http.StatusOK,
			expectedLen:  1,
			assertFn: func(t *testing.T, cars []data.Car) {
				assert.Equal(t, cars[0].Make, "BMW")
			},
		},
		{
			name:         "Filter by model",
			query:        "/v1/cars?model=XC60",
			expectedCode: http.StatusOK,
			expectedLen:  1,
			assertFn: func(t *testing.T, cars []data.Car) {
				assert.Equal(t, cars[0].Model, "XC60")
			},
		},
		{
			name:         "Filter by drivetrain",
			query:        "/v1/cars?drivetrain=AWD",
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "Filter by horsepower range",
			query:        "/v1/cars?horsepower_min=400",
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "Filter by price range",
			query:        "/v1/cars?price_max=55000",
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "No results",
			query:        "/v1/cars?make=Tesla",
			expectedCode: http.StatusOK,
			expectedLen:  0,
		},
		{
			name:         "Pagination works",
			query:        "/v1/cars?page=1&page_size=2",
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name:         "Invalid page",
			query:        "/v1/cars?page=-1",
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "Invalid sort field",
			query:        "/v1/cars?sort=invalid",
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)

			app.listCarsHandler(rr, req)

			assert.Equal(t, rr.Code, tt.expectedCode)

			if tt.expectedCode == http.StatusOK {
				var body struct {
					Cars     []data.Car    `json:"cars"`
					Metadata data.Metadata `json:"metadata"`
				}

				err := json.NewDecoder(rr.Body).Decode(&body)
				if err != nil {
					t.Fatalf("failed to decode JSON: %v", err)
				}

				assert.Equal(t, len(body.Cars), tt.expectedLen)

				if tt.assertFn != nil && len(body.Cars) > 0 {
					tt.assertFn(t, body.Cars)
				}
			}
		})
	}
}
