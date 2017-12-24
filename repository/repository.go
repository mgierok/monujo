package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) (*Repository, error) {
	return &Repository{
		db: db,
	}, nil
}
