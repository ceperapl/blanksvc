package main

import (
	"os"

	"github.com/company/blanksvc/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		os.Exit(1)
	}
}
