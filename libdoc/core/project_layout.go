package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureProjectBaseSymlinks creates genRoot/src and genRoot/config symlinks to
// the product module when those directories exist. credit_framework_core's
// ProjectBasePath walks from the go test cwd (genRoot) looking for both dirs;
// without them, importing framework packages panics: "can not find project bath path".
//
// Symlinks are relative when possible so the gen cache relocates cleanly.
func EnsureProjectBaseSymlinks(genRoot, absModRoot string) error {
	if genRoot == "" || absModRoot == "" {
		return nil
	}
	genRoot, err := filepath.Abs(genRoot)
	if err != nil {
		return err
	}
	absModRoot, err = filepath.Abs(absModRoot)
	if err != nil {
		return err
	}
	for _, name := range []string{"src", "config"} {
		src := filepath.Join(absModRoot, name)
		st, err := os.Stat(src)
		if err != nil || !st.IsDir() {
			continue
		}
		dst := filepath.Join(genRoot, name)
		if _, err := os.Lstat(dst); err == nil {
			// Already present (symlink or dir); leave as-is.
			continue
		}
		rel, err := filepath.Rel(genRoot, src)
		if err != nil {
			rel = src
		}
		if err := os.Symlink(rel, dst); err != nil {
			return fmt.Errorf("symlink %s -> %s: %w", dst, rel, err)
		}
	}
	return nil
}
