package main

import (
	"log"

	"github.com/alianjo/saver/cmd"
)

func main() {
	if err := cmd.NewBackupCMD().Execute(); err != nil {
		log.Fatal(err)
	}
}
