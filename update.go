package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

func GetUpdate(c *cli.Context) {
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
