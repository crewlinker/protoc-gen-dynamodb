package main

import (
	"fmt"
	"os"

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
	return sh.Run("go", "test", "-v", "-failfast", "-count=1", "./...")
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
