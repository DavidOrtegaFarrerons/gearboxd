package data

import (
	"context"
	"database/sql"
	"errors"
	"gearboxd/internal/validator"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrCarLogAlreadyExists = errors.New("a car log for this car and user already exists")
)

type CarLog struct {
	ID        int64            `json:"id"`
	CarID     int64            `json:"-"`
	UserID    int64            `json:"user_id"`
	Rating    *decimal.Decimal `json:"rating"`
	Status    string           `json:"status"`
	Comment   *string          `json:"comment"`
	CreatedAt time.Time        `json:"-"`
	UpdatedAt time.Time        `json:"-"`
}

var (
	statusOptions = []string{"want_to_drive", "driven", "owned"}
)

func ValidateCarLog(v *validator.Validator, cl *CarLog) {
	v.Check(cl.CarID != 0, "car_id", "cannot be 0")
	if cl.Rating != nil {
		v.Check(cl.Rating.GreaterThanOrEqual(decimal.NewFromInt(0)), "rating", "cannot be less than 0")
		v.Check(cl.Rating.LessThanOrEqual(decimal.NewFromInt(5)), "rating", "cannot be more than 5")
	}

	v.Check(validator.PermittedValue(cl.Status, statusOptions...), "status", "status does not exist")
}

type CarLogStore interface {
	Get(id int64) (*CarLog, error)
	Insert(cl *CarLog) error
	Delete(userID, carID int64) error
	GetAllForCar(carID int64) ([]*CarLog, error)
	Update(carLog *CarLog) error
}

type PostgresCarLogStore struct {
	DB *sql.DB
}

func (s *PostgresCarLogStore) Get(id int64) (*CarLog, error) {
	query := `SELECT id, car_id, user_id, rating, status, comment FROM car_logs WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var carLog CarLog
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&carLog.ID,
		&carLog.CarID,
		&carLog.UserID,
		&carLog.Rating,
		&carLog.Status,
		&carLog.Comment,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &carLog, nil
}

func (s *PostgresCarLogStore) Insert(cl *CarLog) error {
	query := `INSERT INTO car_logs(
		user_id, car_id, rating, status, comment
	) VALUES ($1, $2, $3, $4, $5) 
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		cl.UserID,
		cl.CarID,
		cl.Rating,
		cl.Status,
		cl.Comment,
	}

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&cl.ID)
	if err != nil {
		if strings.Contains(err.Error(), "ar_logs_user_id_car_id_key") {
			return ErrCarLogAlreadyExists
		}
	}

	return err
}

func (s *PostgresCarLogStore) Delete(userID, carID int64) error {
	query := `DELETE FROM car_logs WHERE user_id = $1 AND car_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, userID, carID)
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

func (s *PostgresCarLogStore) GetAllForCar(carID int64) ([]*CarLog, error) {
	query := `SELECT id, user_id, rating, status, comment FROM car_logs WHERE car_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query, carID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []*CarLog{}, nil
		default:
			return nil, err
		}
	}

	defer rows.Close()

	var carLogs []*CarLog

	for rows.Next() {
		var carLog CarLog

		err = rows.Scan(
			&carLog.ID,
			&carLog.UserID,
			&carLog.Rating,
			&carLog.Status,
			&carLog.Comment,
		)
		if err != nil {
			return nil, err
		}

		carLogs = append(carLogs, &carLog)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return carLogs, nil
}

func (s *PostgresCarLogStore) Update(carLog *CarLog) error {
	query := `UPDATE car_logs
	SET rating = $1, status = $2, comment = $3 
	WHERE id = $4`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{carLog.Rating, carLog.Status, carLog.Comment, carLog.ID}

	r := s.DB.QueryRowContext(ctx, query, args...)
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}
