package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("no record found")
)

type Models struct {
	Cars CarModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cars: &CarModel{DB: db},
	}
}
