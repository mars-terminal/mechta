package shortener

import (
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	storage *sqlx.DB
}

func NewStorage(storage *sqlx.DB) *Storage {
	return &Storage{storage: storage}
}
