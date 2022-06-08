package main

import (
	"log"
	"os"
	"strconv"

	"github.com/DimitarL/wheeler-trader/server/pkg/app"
)

func main() {
	applicaiton := app.NewApplication()
	host, port := getHostAndPort()

	err := applicaiton.Start(host, port)
	if err != nil {
		log.Fatal(err)
	}
}

func getHostAndPort() (string, int) {
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	if port == "" {
		log.Fatal("APP_PORT environment variable should be provided")
	}
	iPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("APP_PORT has value '%s' which is not valid integer: %s\n", port, err)
	}

	return host, iPort
}
