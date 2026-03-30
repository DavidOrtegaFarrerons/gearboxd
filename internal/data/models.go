package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("no record found")
	ErrEditConflict   = errors.New("conflict found, could not perform operation")
)

type Models struct {
	Cars CarModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cars: &CarModel{DB: db},
	}
}
