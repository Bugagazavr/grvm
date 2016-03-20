package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/codegangsta/cli"
)

var gemfileRegexp = regexp.MustCompile(`^(\s*|)ruby\s*("|')(?P<version>.*)("|')(\s|)$`)

func GetHook(c *cli.Context) {
	var candidate string

	db, err := getDB()
	if err != nil {
		Print(c, "Cannot open database file:", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(false)
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}
	defer tx.Rollback()

	dir, err := os.Getwd()
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	gemfilePath := fmt.Sprintf("%s/%s", dir, "Gemfile")
	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, "Gemfile")); err == nil {
		file, err := os.Open(gemfilePath)
		if err != nil {
			Print(c, err.Error())
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			scanResult := gemfileRegexp.FindStringSubmatch(scanner.Text())
			if len(scanResult) > 3 && len(scanResult[3]) > 0 {
				candidate = scanResult[3]
				break
			}
		}

		if len(candidate) > 0 {
			if err := checkCandidate(tx, candidate); err != nil {
				Print(c, err.Error())
				os.Exit(1)
			}

			printEnv(candidate)
		}
	}

	tx.Commit()
}
