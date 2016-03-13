package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
			Action:  list,
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
			Action:  doctor,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Updates available rubies",
			Action:  update,
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Instqalls ruby",
			Action:  install,
		},
	}

	app.Run(os.Args)
}

func list(c *cli.Context) {
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

func doctor(c *cli.Context) {
	if _, err := os.Stat(rubyBuildDirectory); os.IsNotExist(err) {
		installRubyBuild()
	}

	updateAvailableRubies()
}

func installRubyBuild() {
	fmt.Println("Install ruby-build")
	args := []string{"clone", "https://github.com/rbenv/ruby-build.git", rubyBuildDirectory}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("something going wrong, try to clone ruby-build manually: git", strings.Join(args, " "))
		os.Exit(1)
	}
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

func update(c *cli.Context) {
	if err := os.Chdir(rubyBuildDirectory); err != nil {
		fmt.Println("Cannot switch directory to:", rubyBuildDirectory)
		os.Exit(1)
	}

	args := []string{"pull", "origin", "master"}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("something going wrong, try to update ruby-build manually: git", strings.Join(args, " "))
		os.Exit(1)
	}

	updateAvailableRubies()
}

func updateAvailableRubies() {
	db, err := getDB()
	if err != nil {
		fmt.Println("Cannot open database file:", dbPath)
		os.Exit(1)
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer tx.Rollback()

	b, err := getBucket(tx, []byte("rubies"))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	cmd := exec.Command(rubyBuildExecutable, "--definitions")
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	rubies := strings.Split(string(buffer.Bytes()), "\n")

	for _, ruby := range rubies {
		if len(ruby) != 0 {
			rubyDirectory := fmt.Sprintf("%s/%s", rubiesDirectory, ruby)
			if _, err := os.Stat(rubyDirectory); err == nil {
				b.Put([]byte(ruby), []byte(rubyDirectory))
			} else {
				b.Put([]byte(ruby), make([]byte, 0))
			}
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Cannot commit changes to rubies bucket")
		os.Exit(1)
	}

}

func install(c *cli.Context) {
	installCandidate := c.Args().Get(0)
	candidateDestDirectory := fmt.Sprintf("%s/%s", rubiesDirectory, installCandidate)

	if _, err := os.Stat(candidateDestDirectory); err == nil {
		fmt.Println("You already have installed:", installCandidate)
		os.Exit(1)
	}

	args := []string{installCandidate, candidateDestDirectory}

	cmd := exec.Command(fmt.Sprintf("%s/bin/ruby-build", rubyBuildDirectory), args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Installtion failed")
		os.Exit(1)
	}

	updateAvailableRubies()
}
