package main

import (
	"log"

	"github.com/DimitarL/wheeler-trader/server/pkg/app"
)

func main() {
	applicaiton := app.NewApplication()

	err := applicaiton.Start("localhost", 8080)
	if err != nil {
		log.Fatal(err)
	}
}
