package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

func GetInstall(c *cli.Context) {
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
