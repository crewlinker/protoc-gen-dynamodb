package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/magefile/mage/sh"
)

// Generate generates Go code from definition fiels
func Generate() error {
	if err := sh.Run("buf", "generate",
		"--path", "example/model",
		"--path", "ddb",
	); err != nil {
		return err
	}

	// having the ddb generated file for our own options blocks generation if it fails
	return os.Remove(filepath.Join("proto", "ddb", "v1", "options.ddb.go"))
}

// Checks runs various pre-merge checks
func Checks() error {
	if err := sh.Run("go", "vet", "./..."); err != nil {
		return fmt.Errorf("failed to run go vet: %w", err)
	}

	out, err := sh.Output("go", "fmt", "./...")
	if err != nil {
		return fmt.Errorf("failed to run gofmt: %w", err)
	}

	if out != "" {
		return fmt.Errorf("some files were unformatted, make sure `go fmt` is run")
	}

	if err := sh.Run("go", "run", "-mod=readonly", "honnef.co/go/tools/cmd/staticcheck", "./..."); err != nil {
		return fmt.Errorf("failed to run staticcheck: %w", err)
	}

	if err := sh.Run("buf", "lint"); err != nil {
		return fmt.Errorf("failed to lint protobufs: %w", err)
	}

	return nil
}

// Test test the code base
func Test() error {
	coverdir := "covdatafiles"
	os.RemoveAll(coverdir)
	os.MkdirAll(coverdir, 0777)
	os.Setenv("GOCOVERDIR", coverdir) // allow code generator to write coverage files
	defer os.Unsetenv("GOCOVERDIR")

	if err := Generate(); err != nil {
		return fmt.Errorf("failed to re-generate protobuf code: %w", err)
	}

	if err := sh.Run("go", "run",
		"-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "--fail-on-pending",
		"--junit-report=test-report.xml", "./..."); err != nil {
		return fmt.Errorf("failed to run tests: %w", err)
	}

	if err := sh.Run("go", "tool", "covdata", "textfmt", "-i="+coverdir, "-o="+coverdir+"/out.cover"); err != nil {
		return fmt.Errorf("failed to text format cover data: %w", err)
	}

	return nil
}

// Release tags a new version and pushes it
func Release(version string) error {
	if !regexp.MustCompile(`^v([0-9]+).([0-9]+).([0-9]+)$`).Match([]byte(version)) {
		return fmt.Errorf("version must be in format vX,Y,Z")
	}

	if err := sh.Run("git", "tag", version); err != nil {
		return fmt.Errorf("failed to tag version: %w", err)
	}
	if err := sh.Run("git", "push", "origin", version); err != nil {
		return fmt.Errorf("failed to push version tag: %w", err)
	}

	if err := sh.Run("buf", "push", "-t", version); err != nil {
		return fmt.Errorf("failed to push to buf registry: %w", err)
	}

	return nil
}

// Dev sets up the dev environment using Docker compose
func Dev() error {
	return sh.Run("docker", "compose", "-f", "magefiles/docker-compose.dev.yml", "-p", "protocgenddb-dev", "up",
		"-d", "--build", "--remove-orphans", "--force-recreate")
}

// init performs some sanity checks before running anything
func init() {
	mustBeInRoot()
}

// mustBeInRoot checks that the command is run in the project root
func mustBeInRoot() {
	if _, err := os.Stat("go.mod"); err != nil {
		panic("must be in root, couldn't stat go.mod file: " + err.Error())
	}
}
