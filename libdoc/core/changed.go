package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xhd2015/gitops/git"
)

const NoTestsChangedMessage = "no tests changed"

// ChangedGitFiles resolves the git repository root and returns on-disk changed
// file paths relative to that root.
func ChangedGitFiles(doctestRoot string) (gitRoot string, changedFiles []string, err error) {
	inside, err := git.IsInsideGit(doctestRoot)
	if err != nil {
		return "", nil, err
	}
	if !inside {
		return "", nil, fmt.Errorf("--changed requires a git repository")
	}
	gitRoot, err = git.ShowToplevel(doctestRoot)
	if err != nil {
		return "", nil, fmt.Errorf("--changed requires a git repository")
	}
	changedFiles, err = git.GetOnDiskChangedFiles(gitRoot, git.ResolvePathsToFiles())
	if err != nil {
		return "", nil, err
	}
	return gitRoot, changedFiles, nil
}

// FilterByChangedFiles returns only leaves affected by changed files under doctestRoot.
func FilterByChangedFiles(cases []TreeCase, doctestRoot, gitRoot string, changedFiles []string) []TreeCase {
	if len(cases) == 0 || len(changedFiles) == 0 {
		return nil
	}
	absRoot := canonicalAbsPath(doctestRoot)
	absGitRoot := canonicalAbsPath(gitRoot)

	if doctestRootChanged(absRoot, absGitRoot, changedFiles) {
		return cases
	}

	affected := make(map[string]bool)
	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if !ok {
			continue
		}
		rel = filepath.ToSlash(rel)

		if filepath.Base(rel) == "SETUP.md" {
			groupDir := filepath.ToSlash(filepath.Dir(rel))
			if groupDir == "." {
				groupDir = ""
			}
			for _, tc := range cases {
				if leafAffectedByGroupSetup(tc.Path, groupDir) {
					affected[tc.Path] = true
				}
			}
			continue
		}

		for _, tc := range cases {
			if leafAffectedByPath(tc.Path, rel) {
				affected[tc.Path] = true
			}
		}
	}

	if len(affected) == 0 {
		return nil
	}

	filtered := make([]TreeCase, 0, len(affected))
	for _, tc := range cases {
		if affected[tc.Path] {
			filtered = append(filtered, tc)
		}
	}
	return filtered
}

// ChangedDoctestMarkdownFiles returns absolute paths of changed doctest markdown
// files under doctestRoot.
func ChangedDoctestMarkdownFiles(doctestRoot, gitRoot string, changedFiles []string) []string {
	absRoot := canonicalAbsPath(doctestRoot)
	absGitRoot := canonicalAbsPath(gitRoot)

	var paths []string
	seen := make(map[string]bool)
	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if !ok {
			continue
		}
		base := filepath.Base(rel)
		if base != "DOCTEST.md" && base != "SETUP.md" && base != "ASSERT.md" {
			continue
		}
		absPath := filepath.Join(absRoot, rel)
		if seen[absPath] {
			continue
		}
		seen[absPath] = true
		paths = append(paths, absPath)
	}
	return paths
}

func doctestRootChanged(absRoot, absGitRoot string, changedFiles []string) bool {
	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if ok && rel == "DOCTEST.md" {
			return true
		}
	}
	return false
}

func canonicalAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	canonical, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return abs
	}
	return canonical
}

func relUnder(root, path string) (string, bool) {
	root = canonicalAbsPath(root)
	path = canonicalAbsPath(path)
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return "", false
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return rel, true
}

func leafAffectedByGroupSetup(leafPath, groupDir string) bool {
	if groupDir == "" {
		return true
	}
	if leafPath == groupDir {
		return true
	}
	prefix := groupDir + "/"
	return strings.HasPrefix(leafPath, prefix)
}

func leafAffectedByPath(leafPath, rel string) bool {
	if leafPath == "" {
		return rel == "ASSERT.md" || !strings.Contains(rel, "/")
	}
	if rel == filepath.Join(leafPath, "ASSERT.md") {
		return true
	}
	prefix := leafPath + "/"
	return strings.HasPrefix(rel, prefix)
}