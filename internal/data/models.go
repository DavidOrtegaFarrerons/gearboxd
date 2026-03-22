package data

import "database/sql"

type Models struct {
	Cars CarModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Cars: &CarModel{DB: db},
	}
}
