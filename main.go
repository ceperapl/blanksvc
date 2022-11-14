package main

import (
	"log"

	"github.com/company/blanksvc/cmd"
)

// nolint: lll
//go:generate go-bindata -prefix "pkg/repository/postgres/migrations" -o pkg/repository/postgres/migrations.go -pkg postgres pkg/repository/postgres/migrations/
//go:generate buf generate ./api/protobuf-spec/
//go:generate mockery --dir "./pkg/repository" --filename repository.go --output "./pkg/repository/mocks" --outpkg "mocks" --name Repository
//go:generate mockery --dir "./pkg/service" --filename service.go --output "./pkg/service/mocks" --outpkg "mocks" --name Service

func main() {
	if err := cmd.RunServer(); err != nil {
		log.Fatal(err)
	}
}
