package main

import (
	"github.com/bloznelis/typioca/cmd"
	"log"
)

func main() {
	cmd.OsInit()

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
