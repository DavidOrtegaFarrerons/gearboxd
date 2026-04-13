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
	Cars        CarModelInterface
	Users       UserModelInterface
	Tokens      TokenModelInterface
	Permissions PermissionModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cars:        &CarModel{DB: db},
		Users:       &UserModel{DB: db},
		Tokens:      &TokenModel{DB: db},
		Permissions: &PermissionModel{DB: db},
	}
}
