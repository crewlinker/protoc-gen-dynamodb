package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/magefile/mage/sh"
)

// Generate generates Go code from definition fiels
func Generate() error {
	return sh.Run("buf", "generate")
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

	return nil
}

// Test test the code base
func Test() error {
	return sh.Run("go", "run",
		"-mod=readonly", "github.com/onsi/ginkgo/v2/ginkgo",
		"-p", "-randomize-all", "--fail-on-pending", "--race", "--trace",
		"--junit-report=test-report.xml", "./...")
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
