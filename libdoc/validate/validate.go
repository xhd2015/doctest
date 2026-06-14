package validate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Run(dir string) error {
	return RunWithOptions(dir, core.Options{})
}

func RunWithOptions(dir string, opts core.Options) error {
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

	if opts.Verbose {
		fmt.Printf("[vet] validating %s\n", dir)
	}

	var antiViolations []error

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() && d.Name() == "testdata" {
			return filepath.SkipDir
		}

		if d.IsDir() && path != dir {
			if _, err := os.Stat(filepath.Join(path, "DOCTEST.md")); err == nil {
				if opts.Verbose {
					fmt.Println("[vet] skipping nested DOCTEST.md boundary")
				}
				return filepath.SkipDir
			}
		}

		if d.IsDir() {
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
		}

		name := d.Name()
		if name != "SETUP.md" && name != "ASSERT.md" {
			return nil
		}

		if opts.Verbose {
			fmt.Printf("[vet]   %s\n", name)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		antiViolations = append(antiViolations, checkFileAntiPatterns(path, string(content))...)
		return nil
	})
	if err != nil {
		return err
	}

	return errors.Join(antiViolations...)
}
