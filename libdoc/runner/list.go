package runner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	runnerbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/less-flags"
)

const (
	listANSIGray  = "\x1b[90m"
	listANSIReset = "\x1b[0m"
)

// ListWithWriters inventories doctest roots matching patterns and streams a
// body line per root then a selection-wide summary. Soft empty selection
// returns ErrNoTestsFound (cli soft-exits 0 with "no tests" on stderr).
// Writers are injectable; concurrent calls are independent (no process globals).
func ListWithWriters(args []string, stdout, stderr io.Writer) error {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	_ = stderr // reserved for future warnings; errors return to caller

	colorFlag, noColorFlag, patterns, err := parseListArgs(args)
	if err != nil {
		return err
	}

	mode := core.ColorAuto
	if colorFlag {
		mode = core.ColorAlways
	}
	if noColorFlag {
		mode = core.ColorNever
	}
	// Auto + NO_COLOR non-empty → off (flag --color still wins via ColorAlways).
	if mode == core.ColorAuto && os.Getenv("NO_COLOR") != "" {
		mode = core.ColorNever
	}
	mode = runnerbuild.ResolveColorMode(mode, stdout)
	colorOn := mode == core.ColorAlways

	roots, err := expandListRoots(patterns)
	if err != nil {
		return err
	}

	// Sort by display path for stable, greppable output.
	type rootRow struct {
		abs     string
		display string
	}
	rows := make([]rootRow, 0, len(roots))
	for _, abs := range roots {
		rows = append(rows, rootRow{abs: abs, display: listDisplayPath(abs)})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].display != rows[j].display {
			return rows[i].display < rows[j].display
		}
		return rows[i].abs < rows[j].abs
	})

	var (
		sumLeaves int
		sumL2     int
		sumL3     int
		aggLabels = map[string]int{}
	)

	for _, row := range rows {
		inv, invErr := inventoryRoot(row.abs, row.display)
		if invErr != nil {
			return invErr
		}
		// Stream body immediately (do not buffer all lines then dump).
		fmt.Fprintln(stdout, formatListBodyLine(inv, colorOn))

		sumLeaves += inv.Leaves
		sumL2 += inv.L2
		sumL3 += inv.L3
		for name, n := range inv.Labels {
			aggLabels[name] += n
		}
	}

	// Always include unlabeled in the selection-wide histogram.
	if _, ok := aggLabels["unlabeled"]; !ok {
		aggLabels["unlabeled"] = 0
	}

	// Summary after all body lines: blank, ---, totals, labels.
	fmt.Fprintln(stdout, "")
	sep := "---"
	totals := formatListTotals(len(rows), sumLeaves, sumL2, sumL3)
	labelsLine := "labels: " + formatLabelDist(aggLabels)
	if colorOn {
		sep = listGray(sep)
		totals = listGray(totals)
		labelsLine = listGray(labelsLine)
	}
	fmt.Fprintln(stdout, sep)
	fmt.Fprintln(stdout, totals)
	fmt.Fprintln(stdout, labelsLine)
	return nil
}

func parseListArgs(args []string) (colorFlag, noColorFlag bool, patterns []string, err error) {
	var sawColor, sawNoColor bool
	for _, arg := range args {
		if arg == "--color" {
			sawColor = true
		}
		if arg == "--no-color" {
			sawNoColor = true
		}
	}
	if sawColor && sawNoColor {
		return false, false, nil, fmt.Errorf("--color and --no-color are mutually exclusive")
	}

	remain, err := lessflags.Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
		Parse(args)
	if err != nil {
		return false, false, nil, err
	}
	return colorFlag, noColorFlag, remain, nil
}

// expandListRoots resolves patterns to absolute doctest root paths (deduped).
// Empty patterns default to ./... . Bare "..." is rejected. Empty match returns
// ErrNoTestsFound.
func expandListRoots(patterns []string) ([]string, error) {
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}

	seen := make(map[string]struct{})
	var roots []string
	add := func(p string) error {
		abs, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		abs = filepath.Clean(abs)
		if _, ok := seen[abs]; ok {
			return nil
		}
		seen[abs] = struct{}{}
		roots = append(roots, abs)
		return nil
	}

	for _, arg := range patterns {
		if arg == "..." {
			return nil, fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
		}
		if path_resolve.IsDotDotDotPattern(arg) {
			// List discovers every DOCTEST.md under the base (including nested
			// roots). FindDotDotDotDirs skips children once a parent root is
			// found, which is wrong for inventory of nested trees.
			dirs, err := findAllListRoots(path_resolve.ExtractBasePath(arg))
			if err != nil {
				if errors.Is(err, path_resolve.ErrNoTestsFound) {
					continue
				}
				return nil, err
			}
			for _, d := range dirs {
				if err := add(d); err != nil {
					return nil, err
				}
			}
			continue
		}

		// Plain path: must exist; prefer the path itself when it is a root.
		if _, err := os.Stat(arg); err != nil {
			return nil, err
		}
		abs, err := filepath.Abs(arg)
		if err != nil {
			return nil, err
		}
		abs = filepath.Clean(abs)
		if hasDOCTESTFile(abs) {
			if err := add(abs); err != nil {
				return nil, err
			}
			continue
		}
		if root, ok := path_resolve.ResolveRoot(abs); ok {
			if err := add(root); err != nil {
				return nil, err
			}
			continue
		}
		return nil, fmt.Errorf("no doctest root at %s", arg)
	}

	if len(roots) == 0 {
		return nil, ErrNoTestsFound
	}
	return roots, nil
}

func hasDOCTESTFile(dir string) bool {
	st, err := os.Stat(filepath.Join(dir, "DOCTEST.md"))
	return err == nil && !st.IsDir()
}

// findAllListRoots walks basePath and returns every directory that contains
// DOCTEST.md, including nested roots. Skips testdata/ only. Empty → ErrNoTestsFound.
func findAllListRoots(basePath string) ([]string, error) {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(absBase); err != nil {
		return nil, err
	}
	var dirs []string
	err = filepath.WalkDir(absBase, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if d.Name() == "testdata" {
			return filepath.SkipDir
		}
		if hasDOCTESTFile(path) {
			dirs = append(dirs, path)
			// Continue into children so nested DOCTEST roots are listed too.
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		return nil, ErrNoTestsFound
	}
	return dirs, nil
}

// listDisplayPath prefers a path relative to cwd when under cwd; else absolute.
func listDisplayPath(abs string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return abs
	}
	rel, err := filepath.Rel(cwd, abs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return abs
	}
	if rel == "." {
		return "."
	}
	return rel
}

type rootInventory struct {
	Display string
	Leaves  int
	L2, L3  int
	// Labels includes unlabeled even when 0.
	Labels map[string]int
}

func inventoryRoot(abs, display string) (rootInventory, error) {
	cases, err := core.DiscoverTreeCasesLight(abs)
	if err != nil {
		return rootInventory{}, err
	}
	inv := rootInventory{
		Display: display,
		Leaves:  len(cases),
		Labels:  map[string]int{},
	}
	for _, c := range cases {
		isE2E := false
		if len(c.Labels) == 0 {
			inv.Labels["unlabeled"]++
		} else {
			for _, lab := range c.Labels {
				if lab == "" {
					continue
				}
				inv.Labels[lab]++
				if lab == "e2e" {
					isE2E = true
				}
			}
		}
		if isE2E {
			inv.L3++
		} else {
			inv.L2++
		}
	}
	if _, ok := inv.Labels["unlabeled"]; !ok {
		inv.Labels["unlabeled"] = 0
	}
	return inv, nil
}

func formatListBodyLine(inv rootInventory, colorOn bool) string {
	leavesStr := fmt.Sprintf("%d", inv.Leaves)
	var l2l3 string
	if inv.Leaves == 0 {
		l2l3 = fmt.Sprintf("L2:L3=%d:%d", inv.L2, inv.L3)
	} else {
		p2 := fmt.Sprintf("%.1f", 100*float64(inv.L2)/float64(inv.Leaves))
		p3 := fmt.Sprintf("%.1f", 100*float64(inv.L3)/float64(inv.Leaves))
		l2l3 = fmt.Sprintf("L2:L3=%d:%d (%s%%/%s%%)", inv.L2, inv.L3, p2, p3)
	}
	dist := formatLabelDist(inv.Labels)
	if colorOn {
		leavesStr = listGray(leavesStr)
		l2l3 = listGray(l2l3)
		dist = listGray(dist)
	}
	return inv.Display + "\t" + leavesStr + "\t" + l2l3 + "\t" + dist
}

func formatListTotals(roots, leaves, l2, l3 int) string {
	// Double spaces between major fields (summaryTotalsRE).
	s := fmt.Sprintf("roots=%d  leaves=%d  L2:L3=%d:%d", roots, leaves, l2, l3)
	if leaves > 0 {
		p2 := fmt.Sprintf("%.1f", 100*float64(l2)/float64(leaves))
		p3 := fmt.Sprintf("%.1f", 100*float64(l3)/float64(leaves))
		s += fmt.Sprintf("  (L2 %s%% / L3 %s%%)", p2, p3)
	}
	return s
}

// formatLabelDist returns space-separated name=count, sorted by count desc then name.
func formatLabelDist(m map[string]int) string {
	type kv struct {
		name  string
		count int
	}
	kvs := make([]kv, 0, len(m))
	for name, n := range m {
		kvs = append(kvs, kv{name: name, count: n})
	}
	sort.Slice(kvs, func(i, j int) bool {
		if kvs[i].count != kvs[j].count {
			return kvs[i].count > kvs[j].count
		}
		return kvs[i].name < kvs[j].name
	})
	parts := make([]string, len(kvs))
	for i, k := range kvs {
		parts[i] = fmt.Sprintf("%s=%d", k.name, k.count)
	}
	return strings.Join(parts, " ")
}

func listGray(s string) string {
	return listANSIGray + s + listANSIReset
}
