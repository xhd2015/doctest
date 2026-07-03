package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	modRoot, err := findModuleRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cmd := exec.Command(
		"go", "run", "./script/generate/embed-assert",
		"-o", "libdoc/assertmod/assert.go",
		"-cache-key", "libdoc/assertmod/cache_key.go",
		"-legacy-out", "libdoc/assertmod/legacy_v1",
		"assert",
	)
	cmd.Dir = modRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}