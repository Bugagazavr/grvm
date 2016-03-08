package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Name = "RVM"
	app.Usage = "RVM"
	app.Commands = []cli.Command{
		{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Show env for ruby version",
			Action:  env,
		},
	}

	app.Run(os.Args)
}

func env(c *cli.Context) {
	currentPathEnv := os.Getenv("PATH")
	currentHome := os.Getenv("HOME")
	newPaths := rebuildPaths(currentPathEnv, currentHome)

	currentRuby := "2.3.0"

	fmt.Printf("export GEM_HOME=$HOME/.grvm/gems/%s\n", currentRuby)
	fmt.Printf("export GEM_PATH=$HOME/.grvm/gems/%s\n", currentRuby)
	fmt.Printf("export PATH=%s:$HOME/.grvm/rubies/%s/bin:$HOME/.grvm/gems/%s/bin:$HOME/.grvm/bin\n", newPaths, currentRuby, currentRuby)
}

func rebuildPaths(path, home string) string {
	var paths = strings.Split(path, ":")
	var currentPath = fmt.Sprintf("%s/%s", home, ".grvm")
	var newPaths []string

	for _, p := range paths {
		if !strings.HasPrefix(p, currentPath) {
			newPaths = append(newPaths, p)
		}
	}

	return strings.Join(newPaths, ":")

}
