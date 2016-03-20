package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const (
	rubyBuildRepo    = "https://github.com/rbenv/ruby-build.git"
	latestReleaseUrl = "https://api.github.com/repos/Bugagazavr/grvm/releases/latest"
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
			Usage:   "Installs ruby",
			Action:  GetInstall,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"ui"},
			Usage:   "grvm uninstall <ruby_version>",
			Action:  GetUninstall,
		},
		{
			Name:    "upgrade",
			Aliases: []string{"ug"},
			Usage:   "Get upgrade information",
			Action:  GetUpgrade,
		},
	}

	app.Run(os.Args)
}
