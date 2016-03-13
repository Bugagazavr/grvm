package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

func GetDoctor(c *cli.Context) {
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
