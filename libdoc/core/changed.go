package core

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/gitops/git"
)

const NoTestsChangedMessage = "no tests changed"

// ChangedRunInfo summarizes how --changed filtered a doctest tree.
type ChangedRunInfo struct {
	TotalInTree  int
	ChangedCount int
	Detail       string
}

// ChangedRunInfoForTree computes changed-run metadata for a doctest tree.
func ChangedRunInfoForTree(allCases []TreeCase, doctestRoot, gitRoot string, changedFiles []string) ChangedRunInfo {
	total := len(allCases)
	if len(changedFiles) == 0 {
		return ChangedRunInfo{TotalInTree: total}
	}
	filtered := FilterByChangedFiles(allCases, doctestRoot, gitRoot, changedFiles)
	return ChangedRunInfo{
		TotalInTree:  total,
		ChangedCount: len(filtered),
		Detail:       formatChangedDetail(allCases, filtered, doctestRoot, gitRoot, changedFiles),
	}
}

// ShouldAnnounceChangedRun reports whether to print a --changed status line.
// Zero-change trees are announced only in verbose mode.
func ShouldAnnounceChangedRun(info ChangedRunInfo, verbose bool) bool {
	return verbose || info.ChangedCount > 0
}

// FormatDoctestAnnouncement formats the doctest status line on stderr.
func FormatDoctestAnnouncement(shortPath string, info ChangedRunInfo, changedOnly bool, runningCount int) string {
	if !changedOnly {
		return fmt.Sprintf("doctest: %s (%d tests)", shortPath, runningCount)
	}
	return fmt.Sprintf("doctest: %s (%d tests%s)", shortPath, info.TotalInTree, formatChangedSuffix(info))
}

func formatChangedSuffix(info ChangedRunInfo) string {
	if info.ChangedCount == 0 {
		return ", --changed: 0 tests"
	}
	if info.Detail == "" {
		return fmt.Sprintf(", --changed: %d tests", info.ChangedCount)
	}
	return fmt.Sprintf(", --changed: %d tests, %s", info.ChangedCount, info.Detail)
}

func formatChangedDetail(allCases, filtered []TreeCase, doctestRoot, gitRoot string, changedFiles []string) string {
	if len(filtered) == 0 {
		return ""
	}

	filteredSet := make(map[string]bool, len(filtered))
	for _, tc := range filtered {
		filteredSet[tc.Path] = true
	}

	direct := directlyAffectedLeaves(allCases, doctestRoot, gitRoot, changedFiles)
	assigned := make(map[string]bool, len(direct))
	for leaf := range direct {
		assigned[leaf] = true
	}

	var parts []string
	if n := len(direct); n > 0 {
		parts = append(parts, leafCountPhrase(n))
	}

	absRoot := canonicalAbsPath(doctestRoot)
	absGitRoot := canonicalAbsPath(gitRoot)

	if doctestRootChanged(absRoot, absGitRoot, changedFiles) {
		if others := countUnassigned(filteredSet, assigned); others > 0 {
			parts = append(parts, fmt.Sprintf("1 DOCTEST.md affecting %d other tests", others))
			markUnassigned(filteredSet, assigned)
		}
	}

	if rootSetupChanged(absRoot, absGitRoot, changedFiles) {
		if others := countUnassigned(filteredSet, assigned); others > 0 {
			parts = append(parts, fmt.Sprintf("1 SETUP.md affecting %d other tests", others))
			markUnassigned(filteredSet, assigned)
		}
	}

	for _, groupDir := range changedGroupSetupDirs(absRoot, absGitRoot, changedFiles) {
		others := 0
		for path := range filteredSet {
			if assigned[path] {
				continue
			}
			if leafAffectedByGroupSetup(path, groupDir) {
				others++
				assigned[path] = true
			}
		}
		if others > 0 {
			parts = append(parts, fmt.Sprintf("1 %s/SETUP.md affecting %d other tests", groupDir, others))
		}
	}

	return strings.Join(parts, " + ")
}

func directlyAffectedLeaves(allCases []TreeCase, doctestRoot, gitRoot string, changedFiles []string) map[string]bool {
	affected := make(map[string]bool)
	absRoot := canonicalAbsPath(doctestRoot)
	absGitRoot := canonicalAbsPath(gitRoot)

	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if !ok {
			continue
		}
		rel = filepath.ToSlash(rel)
		if rel == "DOCTEST.md" || filepath.Base(rel) == "SETUP.md" {
			continue
		}
		for _, tc := range allCases {
			if leafAffectedByPath(tc.Path, rel) {
				affected[tc.Path] = true
			}
		}
	}
	return affected
}

func rootSetupChanged(absRoot, absGitRoot string, changedFiles []string) bool {
	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if ok && rel == "SETUP.md" {
			return true
		}
	}
	return false
}

func changedGroupSetupDirs(absRoot, absGitRoot string, changedFiles []string) []string {
	seen := make(map[string]bool)
	var dirs []string
	for _, cf := range changedFiles {
		absChanged := filepath.Join(absGitRoot, cf)
		rel, ok := relUnder(absRoot, absChanged)
		if !ok || filepath.Base(rel) != "SETUP.md" {
			continue
		}
		groupDir := filepath.ToSlash(filepath.Dir(rel))
		if groupDir == "." {
			continue
		}
		if seen[groupDir] {
			continue
		}
		seen[groupDir] = true
		dirs = append(dirs, groupDir)
	}
	sort.Strings(dirs)
	return dirs
}

func countUnassigned(filteredSet, assigned map[string]bool) int {
	n := 0
	for path := range filteredSet {
		if !assigned[path] {
			n++
		}
	}
	return n
}

func markUnassigned(filteredSet, assigned map[string]bool) {
	for path := range filteredSet {
		assigned[path] = true
	}
}

func leafCountPhrase(n int) string {
	if n == 1 {
		return "1 leaf"
	}
	return fmt.Sprintf("%d leaves", n)
}

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
	rel = filepath.ToSlash(rel)
	if leafPath == "" {
		if rel == "ASSERT.md" {
			return true
		}
		return strings.HasPrefix(rel, "testdata/")
	}
	if rel == leafPath+"/ASSERT.md" {
		return true
	}
	return strings.HasPrefix(rel, leafPath+"/testdata/")
}