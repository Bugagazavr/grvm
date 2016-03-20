package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func set(c *cli.Context) {
	db, err := getDB()
	if err != nil {
		Print(c, "Cannot open database file:", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}
	defer tx.Rollback()

	b, err := getBucket(tx, []byte("settings"))
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	candidate := c.Args().Get(0)

	if len(candidate) == 0 {
		Print(c, "No version given")
		os.Exit(1)
	}

	if candidate != "system" {
		if err := checkCandidate(tx, candidate); err != nil {
			Print(c, err.Error())
			os.Exit(1)
		}
	}

	if c.Bool("default") {
		b.Put([]byte("default"), []byte(candidate))

		if err := tx.Commit(); err != nil {
			fmt.Println("Cannot save settings")
			os.Exit(1)
		}

		printEnv(candidate)
		Print(c, "Now,", candidate, "default ruby")
	} else {
		printEnv(candidate)
	}
}
