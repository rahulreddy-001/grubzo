package main

import (
	"log"

	"grubzo/cmd"
)

var (
	version  = "0.1.0"
	revision = "dev"
)

func main() {
	cmd.Version = version
	cmd.Revision = revision
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
