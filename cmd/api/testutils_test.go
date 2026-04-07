package main

import (
	"context"
	"fmt"
	"gearboxd/internal/data"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func newTestApplication(t *testing.T, cfg *config, models *data.Models) *application {
	if cfg == nil {
		cfg = &config{
			env:  "dev",
			port: 4000,
		}
	}

	if models == nil {
		models = &data.Models{
			Cars: &MockCarModel{make([]data.Car, 0)},
		}
	}

	return &application{
		config: *cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: *models,
	}
}

func createTestRequestWithIdParam(t *testing.T, requestMethod, route string, id int, body io.Reader) *http.Request {
	req := httptest.NewRequest(requestMethod, fmt.Sprintf("%s/%d", route, id), body)
	return req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, httprouter.Params{
		{Key: "id", Value: fmt.Sprintf("%d", id)},
	}))
}

type MockCarModel struct {
	cars []data.Car
}

func (m *MockCarModel) Update(car *data.Car) error {
	for i, c := range m.cars {
		if car.ID == c.ID {
			if car.Version != c.Version {
				return data.ErrEditConflict
			}

			car.Version++
			m.cars[i] = *car
			return nil
		}
	}

	return data.ErrRecordNotFound
}

func (m *MockCarModel) Insert(car *data.Car) error {
	car.ID = int64(len(m.cars) + 1)
	m.cars = append(m.cars, *car)
	return nil
}

func (m *MockCarModel) Get(ID int64) (*data.Car, error) {
	for _, car := range m.cars {
		if ID == car.ID {
			return &car, nil
		}
	}

	return nil, data.ErrRecordNotFound
}

func (m *MockCarModel) Delete(id int64) error {
	for i, car := range m.cars {
		if car.ID == id {
			m.cars = append(m.cars[:i], m.cars[i+1:]...)
			return nil
		}
	}

	return data.ErrRecordNotFound
}

func (m *MockCarModel) GetAll(carFilters *data.CarFilters) ([]*data.Car, data.Metadata, error) {
	panic("Implement")
}
