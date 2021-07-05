package repository

import (
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	Db *pgx.Conn
}
