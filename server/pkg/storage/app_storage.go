package storage

import (
	"context"
	"log"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

type AppStorage struct {
	conn *pgx.Conn
}

func NewAppStorage() *AppStorage {
	conn, err := pgx.Connect(context.Background(), getPostgresURI())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &AppStorage{conn: conn}
}

func (a AppStorage) CloseConn() error {
	return a.conn.Close(context.Background())
}

func getPostgresURI() string {
	uri := os.Getenv("APP_POSTGRES_URI")

	if uri == "" {
		log.Fatal("APP_POSTGRES_URI must be provided")
	}

	return uri
}
