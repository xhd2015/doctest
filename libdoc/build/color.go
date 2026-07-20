package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
)

const (
	ansiReset = "\x1b[0m"
	ansiRed   = "\x1b[31m"
	ansiGreen = "\x1b[32m"
	ansiGray  = "\x1b[90m"
)

// isTerminal reports whether w is an *os.File connected to a character device (TTY).
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// ResolveColorMode maps ColorAuto to Always/Never based on whether w is a TTY.
// Explicit ColorAlways and ColorNever are returned unchanged.
//
// Call this against the user-facing stdout (real terminal or pipe) before
// redirecting progress into intermediate buffers (e.g. parallel ./... trees).
// Otherwise ColorAuto type-asserts the buffer as *os.File and always disables color.
func ResolveColorMode(mode core.ColorMode, w io.Writer) core.ColorMode {
	if mode != core.ColorAuto {
		return mode
	}
	if isTerminal(w) {
		return core.ColorAlways
	}
	return core.ColorNever
}

func colorEnabled(mode core.ColorMode, w io.Writer) bool {
	switch ResolveColorMode(mode, w) {
	case core.ColorAlways:
		return true
	default:
		return false
	}
}

type colorStyle struct {
	enabled bool
}

func newColorStyle(mode core.ColorMode, w io.Writer) colorStyle {
	return colorStyle{enabled: colorEnabled(mode, w)}
}

func (c colorStyle) red(s string) string {
	if !c.enabled {
		return s
	}
	return ansiRed + s + ansiReset
}

func (c colorStyle) green(s string) string {
	if !c.enabled {
		return s
	}
	return ansiGreen + s + ansiReset
}

func (c colorStyle) gray(s string) string {
	if !c.enabled {
		return s
	}
	return ansiGray + s + ansiReset
}

// PhaseTiming is a wall-clock span for one pipeline step inside a tree run
// (discover / materialize / generate / go_test / post).
type PhaseTiming struct {
	Name      string // e.g. "discover", "generate", "go_test"
	ElapsedNs int64
}

// LeafTiming is attributed wall time for one leaf when available:
// multi-package trees use go test -json package Elapsed; unified suite trees
// use subtest Elapsed under TestDoctestSuite/<leafPath>. ElapsedNs may be 0
// when unmappable.
type LeafTiming struct {
	Path      string
	ElapsedNs int64
	Cached    bool
}

type TestRunStats struct {
	Passed         int
	Total          int
	Elapsed        time.Duration
	NoTestsChanged bool
	Skipped        []core.SkippedCase
	// Phases are filled by TestWithStats when the tree actually runs.
	Phases []PhaseTiming
	// LeafTimings map leaf paths to go-test package elapsed when available.
	LeafTimings []LeafTiming

	// Filled when GenerateOnly: shared gen module path and tree scope.
	GenRoot string
	TreeRel string
	// Unified is true when hierarchical suite gen was used (not internal-compile).
	Unified bool
	AbsRoot string
}

func formatDisplayDuration(d time.Duration) string {
	if d >= time.Second {
		secs := float64(d) / float64(time.Second)
		s := fmt.Sprintf("%.2f", secs)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s + "s"
	}
	if d >= time.Millisecond {
		return fmt.Sprintf("%dms", d/time.Millisecond)
	}
	if d >= time.Microsecond {
		return fmt.Sprintf("%dµs", d/time.Microsecond)
	}
	return fmt.Sprintf("%dns", d/time.Nanosecond)
}

func formatSummary(style colorStyle, runCount, passCount, failCount, cachedCount int, elapsed time.Duration) string {
	runSeg := fmt.Sprintf("%d Run", runCount)
	passSeg := fmt.Sprintf("%d Pass", passCount)
	failSeg := fmt.Sprintf("%d Fail", failCount)
	cachedSeg := fmt.Sprintf("%d Cached", cachedCount)

	if style.enabled {
		if passCount > 0 {
			passSeg = style.green(passSeg)
		}
		if failCount > 0 {
			failSeg = style.red(failSeg)
		} else {
			failSeg = style.gray(failSeg)
		}
		cachedSeg = style.gray(cachedSeg)
	}

	durSeg := formatDisplayDuration(elapsed)
	if style.enabled {
		durSeg = style.gray(durSeg)
	}

	return fmt.Sprintf("  (%s, %s, %s, %s) in %s", runSeg, passSeg, failSeg, cachedSeg, durSeg)
}

func formatResultSummary(style colorStyle, passed, total int, elapsed time.Duration) string {
	suffix := fmt.Sprintf(" in %s", formatDisplayDuration(elapsed))
	if passed == total {
		token := fmt.Sprintf("PASS (%d/%d)", passed, total)
		if style.enabled {
			return style.green(token) + suffix
		}
		return token + suffix
	}
	token := fmt.Sprintf("FAIL (%d/%d)", passed, total)
	if style.enabled {
		return style.red(token) + suffix
	}
	return token + suffix
}

// PrintSkippedSummary writes a compact skip report grouped by label set.
// Paths and explanations appear only when verbose is true.
func PrintSkippedSummary(skipped []core.SkippedCase, verbose bool) {
	s := FormatSkippedSummary(skipped, verbose)
	if s == "" {
		return
	}
	fmt.Print(s)
}

// FormatSkippedSummary returns the compact skip block (trailing newline included).
//
// Default (verbose=false):
//
//	skipped N labeled (discovery; --label-all or --label EXPR to run)
//	  heavy              12
//	  slow                3
//	  (use -v to list paths)
//
// Verbose: each bucket lists DisplayPath (+ explanation under path).
func FormatSkippedSummary(skipped []core.SkippedCase, verbose bool) string {
	if len(skipped) == 0 {
		return ""
	}

	type bucket struct {
		key   string
		count int
		items []core.SkippedCase
	}
	byKey := map[string]*bucket{}
	for _, s := range skipped {
		key := labelSetKey(s.Labels)
		b := byKey[key]
		if b == nil {
			b = &bucket{key: key}
			byKey[key] = b
		}
		b.count++
		if verbose {
			b.items = append(b.items, s)
		}
	}
	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	// Count desc, then key asc.
	sort.Slice(keys, func(i, j int) bool {
		bi, bj := byKey[keys[i]], byKey[keys[j]]
		if bi.count != bj.count {
			return bi.count > bj.count
		}
		return bi.key < bj.key
	})

	var b strings.Builder
	b.WriteString(skippedSummaryHeader(skipped))
	b.WriteByte('\n')

	// Align counts to a readable column.
	maxKey := 0
	for _, k := range keys {
		if n := len(k); n > maxKey {
			maxKey = n
		}
	}
	if maxKey < 8 {
		maxKey = 8
	}
	if maxKey > 40 {
		maxKey = 40
	}

	for _, k := range keys {
		bk := byKey[k]
		if verbose {
			fmt.Fprintf(&b, "  %s (%d)\n", bk.key, bk.count)
			// Stable path order within bucket.
			items := append([]core.SkippedCase(nil), bk.items...)
			sort.Slice(items, func(i, j int) bool {
				pi, pj := items[i].DisplayPath, items[j].DisplayPath
				if pi == "" {
					pi = items[i].Path
				}
				if pj == "" {
					pj = items[j].Path
				}
				return pi < pj
			})
			for _, it := range items {
				path := it.DisplayPath
				if path == "" {
					path = it.Path
				}
				fmt.Fprintf(&b, "    %s\n", path)
				if it.Explanation != "" {
					fmt.Fprintf(&b, "      explanation: %s\n", it.Explanation)
				}
				if it.Reason != "" {
					fmt.Fprintf(&b, "      reason: %s\n", it.Reason)
				}
			}
			continue
		}
		// Pad key for column alignment (truncate long keys).
		key := bk.key
		if len(key) > maxKey {
			key = key[:maxKey-1] + "…"
		}
		fmt.Fprintf(&b, "  %-*s %d\n", maxKey, key, bk.count)
	}
	if !verbose {
		b.WriteString("  (use -v to list paths)\n")
	}
	b.WriteByte('\n')
	return b.String()
}

func labelSetKey(labels []string) string {
	if len(labels) == 0 {
		return "(unlabeled)"
	}
	cp := append([]string(nil), labels...)
	sort.Strings(cp)
	return strings.Join(cp, ",")
}

func skippedSummaryHeader(skipped []core.SkippedCase) string {
	n := len(skipped)
	allFilter := true
	for _, s := range skipped {
		if s.Reason != "label filter" {
			allFilter = false
			break
		}
	}
	if allFilter {
		return fmt.Sprintf("skipped %d (label filter; --label-all or adjust --label EXPR)", n)
	}
	return fmt.Sprintf("skipped %d labeled (discovery; --label-all or --label EXPR to run)", n)
}

func PrintResultSummary(opts core.Options, stats TestRunStats) {
	if stats.Total == 0 {
		return
	}
	style := newColorStyle(opts.Color, os.Stdout)
	fmt.Println(formatResultSummary(style, stats.Passed, stats.Total, stats.Elapsed))
}

func SkippedDisplayPath(doctestRoot, leafPath string) string {
	short := pathfmt.Short(doctestRoot)
	if leafPath == "" {
		return short
	}
	return short + "/" + filepath.ToSlash(leafPath)
}