package main

import (
	"log"

	"github.com/zam3858/laravel-docgen/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
