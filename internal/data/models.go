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
	Cars        CarStore
	Users       UserStore
	Tokens      PostgresTokenStore
	Permissions PermissionStore
	CarLogs     CarLogStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cars:        &PostgresCarStore{DB: db},
		Users:       &PostgresUserStore{DB: db},
		Tokens:      &TokenModel{DB: db},
		Permissions: &PostgresPermissionStore{DB: db},
		CarLogs:     &PostgresCarLogStore{DB: db},
	}
}
