package main

import (
	"log"

	"github.com/company/blanksvc/cmd"
)

//go:generate go-bindata -prefix "pkg/repository/postgres/migrations" -o pkg/repository/postgres/migrations.go -pkg postgres pkg/repository/postgres/migrations/

func main() {
	if err := cmd.RunServer(); err != nil {
		log.Fatal(err)
	}
}
