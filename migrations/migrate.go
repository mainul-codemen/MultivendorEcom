package main

import (
	"log"

	"github.com/MultivendorEcom/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
