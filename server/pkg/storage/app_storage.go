package storage

import (
	"context"
	"log"

	pgx "github.com/jackc/pgx/v4"
)

type AppStorage struct {
	conn *pgx.Conn
}

func NewAppStorage() *AppStorage {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:0099@localhost:5432/wheeler-trader")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &AppStorage{conn: conn}
}

func (a AppStorage) CloseConn() error {
	return a.conn.Close(context.Background())
}
