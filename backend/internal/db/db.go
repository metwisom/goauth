package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"goauth/internal/config"
)

type db struct {
	pool *pgxpool.Pool
}

func (d *db) Connect() (err error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Config.DbUser, config.Config.DbPassword, config.Config.DbHost, config.Config.DbPort, config.Config.DbName)
	fmt.Println(connStr)

	pgConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	pgConfig.MaxConns = 4 // Adjust this value as needed (default is 4)

	pgConfig.MinConns = 0        // Minimum number of connections to keep open
	pgConfig.MaxConnLifetime = 0 // Maximum lifetime of a connection (0 = no limit)
	pgConfig.MaxConnIdleTime = 0 // Maximum idle time of a connection (0 = no limit)

	d.pool, err = pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return err
	}

	err = d.pool.Ping(context.Background())
	if err != nil {
		return err
	}

	err = d.InitDB()
	return
}

// Close закрывает соединение с базой данных.
func (d *db) Close() error {
	d.pool.Close()
	return nil
}

func (d *db) Exec(sql string, args []any) error {
	_, err := d.pool.Exec(context.Background(), sql, args...)
	return err
}

// Query Выполняет SQL-запрос, возвращающий строки.
func (d *db) Query(sql string, args []any) (*Rows, error) {
	rows, err := d.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return &Rows{Rows: rows}, nil
}

// QueryRow выполняет SQL-запрос, возвращающий одну строку (например, INSERT ... RETURNING).
func (d *db) QueryRow(sql string, args []any, dest ...any) error {
	row := d.pool.QueryRow(context.Background(), sql, args...)
	return row.Scan(dest...)
}

var DB = db{}
