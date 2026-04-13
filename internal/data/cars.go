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
	v.Check(car.Horsepower > 0, "horsepower", "horsepower must be positive and higher than 0")
	v.Check(validator.PermittedValue(car.Fuel, fuelOptions...), "fuel", fmt.Sprintf("fuel must be one of: %+v", fuelOptions))
	v.Check(car.PriceNew.GreaterThan(decimal.NewFromFloat(0)), "price_new", "price_new must be a positive number higher than 0")
}

type CarFilters struct {
	Make          string
	Model         string
	Year          int
	Gearbox       string
	Drivetrain    string
	Fuel          string
	HorsepowerMin int
	HorsepowerMax int
	PriceMin      decimal.Decimal
	PriceMax      decimal.Decimal
	Filters
}

type CarStore interface {
	Insert(car *Car) error
	Get(id int64) (*Car, error)
	Delete(id int64) error
	Update(car *Car) error
	GetAll(carFilters *CarFilters) ([]*Car, Metadata, error)
}
type PostgresCarStore struct {
	DB *sql.DB
}

func (m *PostgresCarStore) Insert(car *Car) error {
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

func (m *PostgresCarStore) Get(id int64) (*Car, error) {
	query := `SELECT id, make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new, version
	FROM cars
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var car Car
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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

func (m *PostgresCarStore) Update(car *Car) error {
	query := `UPDATE cars 
	SET make = $1, model = $2, year = $3, description = $4, image_url = $5, gearbox = $6, drivetrain = $7, horsepower = $8, fuel = $9, price_new = $10, version = version + 1
	WHERE id = $11 AND version = $12
	RETURNING version`

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
		car.ID,
		car.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&car.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m *PostgresCarStore) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM cars WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *PostgresCarStore) GetAll(carFilters *CarFilters) ([]*Car, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, make, model, year, description, image_url, gearbox, drivetrain, horsepower, fuel, price_new, version
	FROM cars
	WHERE
	  (make = $1 OR $1 = '')
	  AND (model ILIKE '%%' || $2 || '%%' OR $2 = '')
	  AND (year = $3 OR $3 = 0)
	  AND (gearbox = $4 OR $4 = '')
	  AND (drivetrain = $5 OR $5 = '')
	  AND (fuel = $6 OR $6 = '')
	  AND (horsepower >= $7 OR $7 = 0)
	  AND (horsepower <= $8 OR $8 = 0)
	  AND (price_new >= $9 OR $9 = 0)
	  AND (price_new <= $10 OR $10 = 0)
	ORDER BY %s %s, id ASC
	LIMIT $11 OFFSET $12
`, carFilters.sortColumn(), carFilters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		carFilters.Make,
		carFilters.Model,
		carFilters.Year,
		carFilters.Gearbox,
		carFilters.Drivetrain,
		carFilters.Fuel,
		carFilters.HorsepowerMin,
		carFilters.HorsepowerMax,
		carFilters.PriceMin,
		carFilters.PriceMax,
		carFilters.Filters.limit(),
		carFilters.Filters.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	cars := []*Car{}

	for rows.Next() {
		var car Car

		err = rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}

		cars = append(cars, &car)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, carFilters.Page, carFilters.PageSize)

	return cars, metadata, nil
}
