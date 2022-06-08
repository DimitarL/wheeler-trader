package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const sqlsDirectory = "file://./sqls"

func main() {
	migrator, err := migrate.New(sqlsDirectory, getPostgresURI())

	if err != nil {
		log.Fatalf("Could not create migrator: %s\n", err)
	}
	defer migrator.Close()

	err = migrator.Migrate(1)
	if err != nil {
		log.Fatalf("Could not migrate db: %s\n", err)
	}

	log.Println("Migration successful")
}

func getPostgresURI() string {
	uri := os.Getenv("APP_POSTGRES_URI")

	if uri == "" {
		log.Fatal("APP_POSTGRES_URI must be provided")
	}

	return uri
}
