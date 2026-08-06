package validate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/rules"
)

// L3 (label: e2e) share budget for full-tree vet. Skipped when ChangedOnly.
const (
	MaxL3Pct  = 10
	MinLeaves = 10
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
	if opts.ChangedOnly {
		return runChangedOnly(dir, opts)
	}
	return runFull(dir, opts)
}

// verboseFprint writes -v progress to opts.Stdout when set (in-process capture),
// otherwise os.Stdout. Parallel-safe: never reassigns process stdout.
func verboseFprint(opts core.Options, format string, args ...any) {
	w := opts.Stdout
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintf(w, format, args...)
}

func verboseFprintln(opts core.Options, a ...any) {
	w := opts.Stdout
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprintln(w, a...)
}

func runChangedOnly(dir string, opts core.Options) error {
	w := opts.Stderr
	if w == nil {
		w = os.Stderr
	}

	gitRoot, changedFiles, err := core.ChangedGitFiles(dir)
	if err != nil {
		return err
	}
	changedPaths := core.ChangedDoctestMarkdownFiles(dir, gitRoot, changedFiles)
	if len(changedPaths) == 0 {
		fmt.Fprintln(w, core.NoTestsChangedMessage)
		return nil
	}

	if opts.Verbose {
		verboseFprint(opts, "[vet] validating %d changed file(s)\n", len(changedPaths))
	}

	var antiViolations []error
	for _, path := range changedPaths {
		if err := validateChangedFile(dir, path, opts, &antiViolations); err != nil {
			return err
		}
	}
	return errors.Join(antiViolations...)
}

func validateChangedFile(dir, path string, opts core.Options, antiViolations *[]error) error {
	name := filepath.Base(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	text := string(content)

	if opts.Verbose {
		rel, _ := filepath.Rel(dir, path)
		verboseFprint(opts, "[vet]   %s\n", rel)
	}

	switch name {
	case "DOCTEST.md":
		if err := checkDOCTESTSections(path, text); err != nil {
			return err
		}
		doctestDoc, err := core.ParseDOCTESTDocument(path, text)
		if err != nil {
			return err
		}
		if doctestDoc.GoBlock == nil {
			return fmt.Errorf("%s: must have a Go code block", path)
		}
		if v := rules.CheckRootHasRequestResponse(doctestDoc.GoBlock.Types, path); v != nil {
			return fmt.Errorf("%s: %s", v.Path, v.Msg)
		}
		if v := rules.CheckRootHasRun(doctestDoc.GoBlock.Run != nil, path); v != nil {
			return fmt.Errorf("%s: %s", v.Path, v.Msg)
		}
		*antiViolations = append(*antiViolations, checkFileAntiPatterns(path, text)...)
	case "SETUP.md":
		if err := checkSETUPSections(path, text); err != nil {
			return err
		}
		setupDoc, parseErr := core.ParseSetupDocument(path, text)
		if parseErr == nil && setupDoc.GoBlock != nil {
			if setupDoc.GoBlock.Setup != nil && filepath.Dir(path) != dir {
				if v := rules.CheckSetupBodyNotStub(setupDoc.GoBlock.Setup.Body, path); v != nil {
					return fmt.Errorf("%s: %s", v.Path, v.Msg)
				}
			}
			if filepath.Dir(path) == dir {
				if v := rules.CheckRootSetupNoRequestResponseRun(setupDoc.GoBlock.Types, setupDoc.GoBlock.Run != nil, path); v != nil {
					return fmt.Errorf("%s: %s", v.Path, v.Msg)
				}
			}
		}
		*antiViolations = append(*antiViolations, checkFileAntiPatterns(path, text)...)
	case "ASSERT.md":
		leafDir := filepath.Dir(path)
		if _, err := os.Stat(filepath.Join(leafDir, "SETUP.md")); err != nil {
			return fmt.Errorf("%s: ASSERT.md found but SETUP.md missing", leafDir)
		}
		if _, err := core.ParseAssertDocument(path, text); err != nil {
			return err
		}
		*antiViolations = append(*antiViolations, checkFileAntiPatterns(path, text)...)
	}
	return nil
}

func runFull(dir string, opts core.Options) error {
	doctestPath := filepath.Join(dir, "DOCTEST.md")
	if _, err := os.Stat(doctestPath); err != nil {
		return fmt.Errorf("%s: root must contain DOCTEST.md", dir)
	}
	doctestContent, err := os.ReadFile(doctestPath)
	if err != nil {
		return err
	}
	if err := checkDOCTESTSections(doctestPath, string(doctestContent)); err != nil {
		return err
	}
	doctestDoc, err := core.ParseDOCTESTDocument(doctestPath, string(doctestContent))
	if err != nil {
		return err
	}
	if doctestDoc.GoBlock == nil {
		return fmt.Errorf("%s: must have a Go code block", doctestPath)
	}
	if v := rules.CheckRootHasRequestResponse(doctestDoc.GoBlock.Types, doctestPath); v != nil {
		return fmt.Errorf("%s: %s", v.Path, v.Msg)
	}
	if v := rules.CheckRootHasRun(doctestDoc.GoBlock.Run != nil, doctestPath); v != nil {
		return fmt.Errorf("%s: %s", v.Path, v.Msg)
	}
	antiViolations := checkFileAntiPatterns(doctestPath, string(doctestContent))

	if opts.Verbose {
		verboseFprint(opts, "[vet] validating %s\n", dir)
	}

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
					verboseFprintln(opts, "[vet] skipping nested DOCTEST.md boundary")
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
			verboseFprint(opts, "[vet]   %s\n", name)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(content)
		if name == "ASSERT.md" {
			if _, err := core.ParseAssertDocument(path, text); err != nil {
				return err
			}
		}
		if name == "SETUP.md" {
			if err := checkSETUPSections(path, text); err != nil {
				return err
			}
			setupDoc, parseErr := core.ParseSetupDocument(path, text)
			if parseErr == nil && setupDoc.GoBlock != nil {
				if setupDoc.GoBlock.Setup != nil && filepath.Dir(path) != dir {
					if v := rules.CheckSetupBodyNotStub(setupDoc.GoBlock.Setup.Body, path); v != nil {
						return fmt.Errorf("%s: %s", v.Path, v.Msg)
					}
				}
				if filepath.Dir(path) == dir {
					if v := rules.CheckRootSetupNoRequestResponseRun(setupDoc.GoBlock.Types, setupDoc.GoBlock.Run != nil, path); v != nil {
						return fmt.Errorf("%s: %s", v.Path, v.Msg)
					}
				}
			}
		}
		antiViolations = append(antiViolations, checkFileAntiPatterns(path, text)...)
		return nil
	})
	if err != nil {
		return err
	}

	// Full tree only: enforce L3 (e2e) share budget after structure + anti-patterns.
	// Skipped entirely when opts.ChangedOnly (runChangedOnly path never reaches here).
	if shareErr := checkL3ShareBudget(dir); shareErr != nil {
		antiViolations = append(antiViolations, shareErr)
	}

	return errors.Join(antiViolations...)
}

// checkL3ShareBudget fails when leaves >= MinLeaves and e2e share exceeds MaxL3Pct.
// L3 identity matches doctest list: only label "e2e" (non-e2e labels are L2).
func checkL3ShareBudget(root string) error {
	cases, err := core.DiscoverTreeCasesLight(root)
	if err != nil {
		return err
	}
	leaves := len(cases)
	if leaves < MinLeaves {
		return nil
	}
	l3 := 0
	for _, c := range cases {
		for _, lab := range c.Labels {
			if lab == "e2e" {
				l3++
				break
			}
		}
	}
	pct := 100.0 * float64(l3) / float64(leaves)
	if pct <= float64(MaxL3Pct) {
		return nil
	}
	return fmt.Errorf("%s: L3 share %.1f%% (%d/%d leaves labeled e2e) exceeds max %d%%",
		root, pct, l3, leaves, MaxL3Pct)
}
