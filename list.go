package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func GetList(c *cli.Context) {
	var installedRubies []string
	var knownRubies []string

	db, err := getDB()
	if err != nil {
		Print(c, "Cannot open database file:", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("rubies"))
	cursor := b.Cursor()

	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		if v != nil && len(v) != 0 {
			installedRubies = append(installedRubies, string(k))
		} else {
			knownRubies = append(knownRubies, string(k))
		}
	}

	fmt.Println("installed rubies:")
	for _, ruby := range installedRubies {
		fmt.Println("  ", ruby)
	}

	if c.Bool("known") {
		fmt.Println("known rubies:")
		for _, ruby := range knownRubies {
			fmt.Println("  ", ruby)
		}
	}
}
