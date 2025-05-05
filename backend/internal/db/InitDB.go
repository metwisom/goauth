package db

import (
	"context"
)

// InitDB Инициирует соединение с БД
func (d *db) InitDB() (err error) {
	_, err = d.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        login text UNIQUE,
        password TEXT,
        username TEXT,
        steam_id TEXT
    )`)
	if err != nil {
		return
	}

	// Создание таблицы сессий
	_, err = d.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS sessions (
        session_id TEXT PRIMARY KEY,
        user_id INTEGER,
        expires_at TIMESTAMP
    )`)
	if err != nil {
		return
	}

	// Создание таблицы одноразовых кодов
	_, err = d.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS codes (
        code TEXT PRIMARY KEY,
        client_id INTEGER,
        user_id INTEGER,
        expires_at TIMESTAMP
    )`)
	if err != nil {
		return
	}

	_, err = d.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS access_token (
        "access_token" VARCHAR NOT NULL,
        "user_id" INTEGER NOT NULL,
        "expires_at" TIMESTAMP NOT NULL
    )`)
	if err != nil {
		return
	}

	// Создание последовательности для client_client_id_seq
	_, err = d.pool.Exec(context.Background(), `CREATE SEQUENCE IF NOT EXISTS client_client_id_seq
        MINVALUE 1
        MAXVALUE 2147483647
        START 1
        CACHE 1`)
	if err != nil {
		return
	}

	// Создание таблицы client
	_, err = d.pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS client (
        "user_id" INTEGER NOT NULL,
        "secret" VARCHAR NOT NULL,
        "client_id" BIGINT DEFAULT nextval('client_client_id_seq'::regclass) NOT NULL
    )`)
	if err != nil {
		return
	}

	// Создание индекса на user_id
	_, err = d.pool.Exec(context.Background(), `CREATE INDEX IF NOT EXISTS index_user_id ON client (user_id)`)
	if err != nil {
		return
	}
	return
}
