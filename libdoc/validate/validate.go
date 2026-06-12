package validate

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}
	if _, err := os.Stat(filepath.Join(dir, "DOCTEST.md")); err != nil {
		return fmt.Errorf("%s: root must contain DOCTEST.md", dir)
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() || d.Name() == "testdata" {
			if d.IsDir() && d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		hasSetup := false
		hasAssert := false
		for _, ent := range entries {
			switch ent.Name() {
			case "SETUP.md":
				hasSetup = true
			case "ASSERT.md":
				hasAssert = true
			}
		}
		if hasAssert && !hasSetup {
			return fmt.Errorf("%s: ASSERT.md found but SETUP.md missing", path)
		}
		return nil
	})
}
