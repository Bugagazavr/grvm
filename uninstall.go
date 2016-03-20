package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

func GetUninstall(c *cli.Context) {
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

	candidate := c.Args().Get(0)

	if len(candidate) == 0 {
		Print(c, "No version given")
		os.Exit(1)
	}

	if err := checkCandidate(tx, candidate); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	if err := exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", rubiesDirectory, candidate)).Run(); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	if err := exec.Command("rm", "-rf", fmt.Sprintf("%s/%s", gemsDirectory, candidate)).Run(); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	if err := updateAvailableRubiesWithTx(tx); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	if err := tx.Commit(); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
