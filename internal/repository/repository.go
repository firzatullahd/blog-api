package repository

import (
	"database/sql"
)

type Repo struct {
	db     *sql.DB
	dbRead *sql.DB
}

func NewRepository(master, replica *sql.DB) *Repo {
	return &Repo{
		db:     master,
		dbRead: replica,
	}
}

func (r *Repo) WithTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}
