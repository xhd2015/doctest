package runner

import (
	"os"
	"path/filepath"
)

// resolveTestTarget maps a CLI path to the doctest subdir and whether it names a concrete leaf.
func resolveTestTarget(arg string) (targetDir string, explicitLeaf bool) {
	abs, err := filepath.Abs(arg)
	if err != nil {
		return arg, false
	}
	if filepath.Base(abs) == "ASSERT.md" {
		return filepath.Dir(abs), true
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return abs, false
	}
	if _, err := os.Stat(filepath.Join(abs, "ASSERT.md")); err == nil {
		return abs, true
	}
	return abs, false
}