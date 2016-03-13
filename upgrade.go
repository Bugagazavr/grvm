package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
)

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

func GetUpgrade(c *cli.Context) {
	resp, err := http.Get(latestReleaseUrl)
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		Print(c, err.Error())
		os.Exit(1)
	}

	var releaseVersion = strings.TrimPrefix(release.TagName, "v")

	if releaseVersion == version {
		Print(c, "You have already latest version:", releaseVersion)
		if c.GlobalBool("shell") {
			fmt.Println("unset grvm_upgrade_url")
			fmt.Println("unset grvm_upgrade_version")
		}
		os.Exit(0)
	}

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, runtime.GOARCH) && strings.Contains(asset.Name, runtime.GOOS) {
			if c.GlobalBool("shell") {
				fmt.Println(Export("grvm_upgrade_url", asset.DownloadUrl))
				fmt.Println(Export("grvm_upgrade_version", releaseVersion))
				break
			} else {
				fmt.Println("New version available:", asset.DownloadUrl)
			}
		}
	}
}
