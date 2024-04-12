package main

import (
	"context"
	"log"
	"os"

	"github.com/pressly/goose/v3"

	_ "cms/migrations"

	_ "github.com/lib/pq"
)

//nolint:gocritic
func main() {
	database, err := goose.OpenDBWithDriver(os.Getenv("GOOSE_DRIVER"), os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if err := goose.RunContext(context.Background(), os.Args[1], database, os.Args[2]); err != nil {
		log.Fatalf("goose %v: %v", os.Args[1:], err)
	}
}
