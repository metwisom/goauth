package db

import (
	"github.com/jackc/pgx/v5"
)

type Rows struct {
	Rows pgx.Rows
}

// Get извлекает следующую строку из результата запроса.
func (rows *Rows) Get(dest ...any) bool {
	return rows.Get(dest...)
}

// Close закрывает результат запроса.
func (rows *Rows) Close() error {
	return rows.Close()
}
