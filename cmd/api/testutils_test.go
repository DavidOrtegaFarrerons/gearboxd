package main

import (
	"gearboxd/internal/data"
	"io"
	"log/slog"
	"testing"
)

func newTestApplication(t *testing.T, cfg *config) *application {
	if cfg == nil {
		cfg = &config{
			env:  "dev",
			port: 4000,
		}
	}

	return &application{
		config: *cfg,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: data.Models{
			Cars: &MockCarModel{make([]data.Car, 0)},
		},
	}
}

type MockCarModel struct {
	cars []data.Car
}

func (m *MockCarModel) Insert(car *data.Car) error {
	m.cars = append(m.cars, *car)
	return nil
}
