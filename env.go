package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func GetEnv(c *cli.Context) {
	var currentRuby string
	var err error

	if len(grvmRubyEnv) == 0 {
		currentRuby, err = getDefaultRuby(c)
		if err != nil {
			Print(c, "Cannot get default ruby:", err.Error())
		}
	} else {
		currentRuby = grvmRubyEnv
	}

	switch currentRuby {
	case "system":
		fmt.Println("export grvm_ruby=system")
	default:
		printEnv(currentRuby)
	}
}
