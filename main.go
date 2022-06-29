package main

import (
	"log"

	"github.com/bloznelis/typioca/cmd"
)

func main() {
	cmd.OsInit()

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
