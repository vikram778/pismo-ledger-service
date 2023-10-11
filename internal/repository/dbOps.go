package repository

import "github.com/jmoiron/sqlx"

type dBOps struct {
	db *sqlx.DB
}

// NewDBOpsRepository initialises db operation repo
func NewDBOpsRepository(db *sqlx.DB) dBOps {
	return dBOps{db: db}
}
