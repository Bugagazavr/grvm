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

func getDB() (*bolt.DB, error) {
	return bolt.Open(dbPath, 0600, nil)
}

func Print(c *cli.Context, args ...string) {
	if c.GlobalBool("shell") {
		fmt.Println(fmt.Sprintf("echo %q", strings.Join(args, " ")))
	} else {
		fmt.Println(strings.Join(args, " "))
	}
}

func Export(key, value string) string {
	return fmt.Sprintf("export %s=%s", key, value)
}

func printEnv(rubyVersion string) {
	newPaths := rebuildPaths()

	if rubyVersion == "system" {
		fmt.Println(Export("PATH", newPaths))
	} else {
		gemsPath := fmt.Sprintf("%s/%s", gemsDirectory, rubyVersion)

		fmt.Println(Export("GEM_HOME", gemsPath))
		fmt.Println(Export("GEM_PATH", gemsPath))

		currentRubyBin := fmt.Sprintf("%s/%s/bin", rubiesDirectory, rubyVersion)
		currentGemsBin := fmt.Sprintf("%s/bin", gemsPath)

		path := fmt.Sprintf("%s:%s:%s", currentRubyBin, currentGemsBin, newPaths)
		fmt.Println(Export("PATH", path))
	}

	fmt.Println(Export("grvm_ruby", rubyVersion))
}

func rebuildPaths() string {
	paths := strings.Split(currentPathEnv, ":")

	var currentPath = fmt.Sprintf("%s/%s", currentHomeEnv, ".grvm")
	var newPaths []string

	for _, p := range paths {
		if !strings.HasPrefix(p, currentPath) {
			newPaths = append(newPaths, p)
		}
	}

	return strings.Join(newPaths, ":")
}

func updateRubiesBucket(b *bolt.Bucket) error {
	buffer := bytes.NewBuffer(make([]byte, 0))

	cmd := exec.Command(rubyBuildExecutable, "--definitions")
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	if err := cmd.Run(); err != nil {
		return err
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

	return nil
}

func updateAvailableRubiesWithTx(tx *bolt.Tx) error {
	b, err := getBucket(tx, []byte("rubies"))
	if err != nil {
		return err
	}

	return updateRubiesBucket(b)
}

func updateAvailableRubies() error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateAvailableRubiesWithTx(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func getDefaultRuby(c *cli.Context) (string, error) {
	db, err := getDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	tx, err := db.Begin(true)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	b, err := getBucket(tx, []byte("settings"))
	if err != nil {
		return "", err
	}

	defaultRuby := b.Get([]byte("default"))

	if defaultRuby == nil {
		return "system", nil
	} else if len(defaultRuby) == 0 {
		return "system", nil
	} else {
		candidate := string(defaultRuby)

		if err := checkCandidate(tx, candidate); err != nil {
			return "system", nil
		} else {
			return candidate, nil
		}
	}
}

func getBucket(tx *bolt.Tx, bucket []byte) (*bolt.Bucket, error) {
	var b *bolt.Bucket
	var e error

	b, e = tx.CreateBucket(bucket)

	if e == bolt.ErrBucketExists {
		return tx.Bucket(bucket), nil
	} else if e == nil {
		return b, e
	} else {
		return nil, e
	}
}

func checkCandidate(tx *bolt.Tx, candidate string) error {
	rubies := tx.Bucket([]byte("rubies"))
	value := rubies.Get([]byte(candidate))

	if value == nil {
		return fmt.Errorf("No candidate to set: %s", candidate)
	} else if len(value) == 0 {
		return fmt.Errorf("%s not installed, please use: grvm install %s", candidate, candidate)
	}

	return nil
}
