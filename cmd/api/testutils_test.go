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
	"time"

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

func (m *MockCarModel) GetAll(cf *data.CarFilters) ([]*data.Car, data.Metadata, error) {
	var filtered []*data.Car

	for i := range m.cars {
		car := m.cars[i]

		if cf.Make != "" && car.Make != cf.Make {
			continue
		}

		if cf.Model != "" && car.Model != cf.Model {
			continue
		}

		if cf.Year != 0 && car.Year != cf.Year {
			continue
		}

		if cf.Gearbox != "" && car.Gearbox != cf.Gearbox {
			continue
		}

		if cf.Drivetrain != "" && car.Drivetrain != cf.Drivetrain {
			continue
		}

		if cf.Fuel != "" && car.Fuel != cf.Fuel {
			continue
		}

		if cf.HorsepowerMin != 0 && car.Horsepower < cf.HorsepowerMin {
			continue
		}

		if cf.HorsepowerMax != 0 && car.Horsepower > cf.HorsepowerMax {
			continue
		}

		if !cf.PriceMin.IsZero() && car.PriceNew.LessThan(cf.PriceMin) {
			continue
		}

		if !cf.PriceMax.IsZero() && car.PriceNew.GreaterThan(cf.PriceMax) {
			continue
		}

		filtered = append(filtered, &car)
	}

	total := len(filtered)

	start := (cf.Filters.Page - 1) * cf.Filters.PageSize
	end := start + cf.Filters.PageSize

	if start > total {
		filtered = []*data.Car{}
	} else {
		if end > total {
			end = total
		}
		filtered = filtered[start:end]
	}

	metadata := data.Metadata{
		CurrentPage:  cf.Filters.Page,
		PageSize:     cf.Filters.PageSize,
		TotalRecords: total,
	}

	return filtered, metadata, nil
}

type MockUserStore struct {
	InsertFunc      func(user *data.User) error
	UpdateFunc      func(user *data.User) error
	GetByEmailFunc  func(email string) (*data.User, error)
	GetForTokenFunc func(scope, tokenPlaintext string) (*data.User, error)
}

func (m *MockUserStore) Insert(user *data.User) error {
	return m.InsertFunc(user)
}

func (m *MockUserStore) Update(user *data.User) error {
	return m.UpdateFunc(user)
}

func (m *MockUserStore) GetByEmail(email string) (*data.User, error) {
	return m.GetByEmailFunc(email)
}

func (m *MockUserStore) GetForToken(scope, tokenPlaintext string) (*data.User, error) {
	return m.GetForTokenFunc(scope, tokenPlaintext)
}

type MockTokenStore struct {
	NewFunc              func(userID int64, ttl time.Duration, scope string) (*data.Token, error)
	InsertFunc           func(token *data.Token) error
	DeleteAllForUserFunc func(scope string, userID int64) error
}

func (m MockTokenStore) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	return m.NewFunc(userID, ttl, scope)
}

func (m MockTokenStore) Insert(token *data.Token) error {
	return m.InsertFunc(token)
}

func (m MockTokenStore) DeleteAllForUser(scope string, userID int64) error {
	return m.DeleteAllForUserFunc(scope, userID)
}

type MockPermissionStore struct {
	GetAllForUserFunc func(userID int64) (data.Permissions, error)
	AddForUserFunc    func(userID int64, codes ...string) error
}

func (m MockPermissionStore) GetAllForUser(userID int64) (data.Permissions, error) {
	return m.GetAllForUserFunc(userID)
}

func (m MockPermissionStore) AddForUser(userID int64, codes ...string) error {
	return m.AddForUserFunc(userID, codes...)
}
