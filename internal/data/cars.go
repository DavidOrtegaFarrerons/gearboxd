package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gearboxd/internal/validator"
	"time"

	"github.com/shopspring/decimal"
)

type Car struct {
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
	CreatedAt   time.Time       `json:"-"`
	UpdatedAt   time.Time       `json:"-"`
	DeletedAt   *time.Time      `json:"-"`
}

var (
	gearboxOptions    = []string{"manual", "automatic", "DCT", "CVT"}
	drivetrainOptions = []string{"FWD", "RWD", "AWD", "4WD"}
	fuelOptions       = []string{"diesel", "gas", "electric", "hybrid", "plug-in-hybrid", "hydrogen", "lpg", "cng"}
)

func ValidateCar(v *validator.Validator, car *Car) {
	v.Check(car.Make != "", "make", "make cannot be empty")
	v.Check(car.Model != "", "model", "model cannot be empty")
	v.Check(car.Year != 0, "year", "year cannot be empty")
	v.Check(car.Year > 1886, "year", "the minimum year is 1886, when the first car was patented")
	v.Check(car.Description != "", "description", "description cannot be empty")
	v.Check(car.ImageURL != "", "image_url", "image_url cannot be empty")
	v.Check(validator.PermittedValue(car.Gearbox, gearboxOptions...), "gearbox", fmt.Sprintf("gearbox must be one of: %+v", gearboxOptions))
	v.Check(validator.PermittedValue(car.Drivetrain, drivetrainOptions...), "drivetrain", fmt.Sprintf("drivetrain must be one of: %+v", drivetrainOptions))
	v.Check(validator.PermittedValue(car.Fuel, fuelOptions...), "fuel", fmt.Sprintf("fuel must be one of: %+v", fuelOptions))
	v.Check(car.PriceNew.GreaterThan(decimal.NewFromFloat(0)), "price_new", "price_new must be a positive number higher than 0")
}

type CarModelInterface interface {
	Insert(car *Car) error
	Get(ID int64) (*Car, error)
}
type CarModel struct {
	DB *sql.DB
}

func (m *CarModel) Insert(car *Car) error {
	query := `INSERT INTO cars (
	  make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, created_at, version
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		car.Make,
		car.Model,
		car.Year,
		car.Description,
		car.ImageURL,
		car.Gearbox,
		car.Drivetrain,
		car.Horsepower,
		car.Fuel,
		car.PriceNew,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&car.ID, &car.CreatedAt, &car.Version)
}

func (m *CarModel) Get(ID int64) (*Car, error) {
	query := `SELECT id, make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new, version
	FROM cars
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var car Car
	err := m.DB.QueryRowContext(ctx, query, ID).Scan(
		&car.ID,
		&car.Make,
		&car.Model,
		&car.Year,
		&car.Description,
		&car.ImageURL,
		&car.Gearbox,
		&car.Drivetrain,
		&car.Horsepower,
		&car.Fuel,
		&car.PriceNew,
		&car.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &car, nil
}
