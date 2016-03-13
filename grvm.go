package main

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/cli"
)

const (
	rubyBuildRepo = "https://github.com/rbenv/ruby-build.git"
)

// System
var version = ""

// ENVs
var currentHomeEnv = os.Getenv("HOME")
var grvmRubyEnv = os.Getenv("grvm_ruby")
var currentPathEnv = os.Getenv("PATH")

// Directories
var grvmDirectory = fmt.Sprintf("%s/.grvm", currentHomeEnv)
var rubyBuildDirectory = fmt.Sprintf("%s/ruby-build", grvmDirectory)
var rubiesDirectory = fmt.Sprintf("%s/rubies", grvmDirectory)
var gemsDirectory = fmt.Sprintf("%s/gems", grvmDirectory)

// Paths
var rubyBuildExecutable = fmt.Sprintf("%s/bin/ruby-build", rubyBuildDirectory)
var dbPath = fmt.Sprintf("%s/grvm.db", grvmDirectory)

func main() {
	app := cli.NewApp()
	app.Name = "GRVM"
	app.Usage = "GRVM"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "shell, s",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Show list installed rubies",
			Action:  GetList,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "known, k",
				},
			},
		},
		{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Show env for ruby version",
			Action:  GetEnv,
		},
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "Set current ruby",
			Action:  set,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "default, d",
				},
			},
		},
		{
			Name:    "doctor",
			Aliases: []string{"d"},
			Usage:   "Fixes all bugs",
			Action:  GetDoctor,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Updates available rubies",
			Action:  GetUpdate,
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Instqalls ruby",
			Action:  GetInstall,
		},
	}

	app.Run(os.Args)
}

func set(c *cli.Context) {
	db, err := getDB()
	if err != nil {
		Print(c, "Cannot open database file:", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tx.Rollback()

	var b *bolt.Bucket
	var e error
	b, e = tx.CreateBucket([]byte("settings"))
	if e == bolt.ErrBucketExists {
		b = tx.Bucket([]byte("settings"))
	} else if e != nil {
		Print(c, "Cannot create bucket for settings")
		os.Exit(1)
	}

	candidate := c.Args().Get(0)

	if len(candidate) == 0 {
		Print(c, "No version given")
		os.Exit(1)
	}

	if err := checkCandidate(tx, candidate); err != nil {
		Print(c, err.Error())
		os.Exit(1)
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
