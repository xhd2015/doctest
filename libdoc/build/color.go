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
	ansiReset  = "\x1b[0m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiGray   = "\x1b[90m"
	ansiOrange = "\x1b[38;5;208m" // warning accent (timeout cancelled segment)
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

func (c colorStyle) orange(s string) string {
	if !c.enabled {
		return s
	}
	return ansiOrange + s + ansiReset
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
	// SkipCount is suite-leaf go test Action "skip" (runtime t.Skip). Distinct
	// from Skipped (label-discovery skips shown in the discovery block).
	SkipCount int
	// Planned is discovery leaf count before Total is rewritten to actual_run
	// (pass+fail). Used for timeout cancelled accounting.
	Planned int
	// TimedOut is true when this run surfaced a go test timeout Error
	// ("test timed out after …"). Enables FAIL (p/planned, N cancelled).
	TimedOut bool
	// BuildFailed is true when go test reported [build failed] with no suite
	// leaf events. Summary must not invent Run/Cached from package fail +
	// pre-planned leaf-cache skips; footer uses planned denom.
	BuildFailed bool
	// Phases are filled by TestWithStats when the tree actually runs.
	Phases []PhaseTiming
	// LeafTimings map leaf paths to go-test package elapsed when available.
	LeafTimings []LeafTiming

	// Filled when GenerateOnly: shared gen module path and tree scope.
	GenRoot string
	TreeRel string
	// SuiteRel is gen-relative suite placement (path-local when PathScoped).
	SuiteRel string
	// PathScoped is true when the run is under a mid/leaf path, not whole tree.
	PathScoped bool
	// Unified is true when hierarchical suite gen was used (not internal-compile).
	Unified bool
	AbsRoot string

	// GoTestBypassed is true when DOCTEST_DEBUG bypass-go-test skipped go test
	// exec after successful prepare (and workspace write when applicable).
	GoTestBypassed bool

	// Cases are runnable leaves after label/changed filters (set by TestWithStats).
	// PrepareTree reuses this to avoid a second DiscoverTreeCases walk.
	Cases []core.TreeCase
}

// timeoutCancelled returns cancelled leaf count when TimedOut; else 0.
// cancelled = max(0, planned − pass − fail − skipCount) with fail derived
// from actual_run Total (pass+fail) when Total is the post-run actual_run.
func timeoutCancelled(stats TestRunStats) int {
	if !stats.TimedOut {
		return 0
	}
	planned := stats.Planned
	if planned <= 0 {
		// Fallback: Total may still hold discovery count when no actual_run.
		planned = stats.Total
	}
	fail := stats.Total - stats.Passed
	if fail < 0 {
		fail = 0
	}
	// When Total still equals planned (no actual_run rewrite), treat unfinished
	// as cancelled rather than as fails: use fail only from actual finished fails.
	// actual_run = pass+fail ≤ planned; if Total == planned and TimedOut with
	// no per-leaf results, Passed was forced to 0 and Total may be 0 (see apply).
	c := planned - stats.Passed - fail - stats.SkipCount
	if c < 0 {
		return 0
	}
	return c
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

// formatResultSummary builds the end-of-run PASS/FAIL line.
// When forceFail is true (e.g. a sibling tree failed prepare while survivors
// all passed), always print FAIL so the summary matches non-zero exit.
//
// When cancelled > 0 (go test timeout with unfinished planned leaves), fraction
// is passed/planned with ", N cancelled" (v1: omit t.Skip on this line). The
// cancelled phrase uses orange when color is on; FAIL token stays red.
//
// When skipCount > 0 (runtime t.Skip on suite leaves) and no timeout cancelled,
// fraction is succeeded/actual_run with ", N t.Skip" (actual_run = pass+fail,
// excludes skips). When skipCount == 0 and cancelled == 0, keep PASS (p/t) /
// FAIL (p/t) unchanged.
func formatResultSummary(style colorStyle, passed, total int, elapsed time.Duration, forceFail bool, skipCount int, cancelled int) string {
	suffix := fmt.Sprintf(" in %s", formatDisplayDuration(elapsed))
	if cancelled > 0 {
		// Timeout path: always FAIL with planned denom + cancelled segment.
		cancelledSeg := fmt.Sprintf("%d cancelled", cancelled)
		if style.enabled {
			// Prefer orange only on the cancelled phrase; FAIL token red.
			return style.red("FAIL") + fmt.Sprintf(" (%d/%d, %s)", passed, total, style.orange(cancelledSeg)) + suffix
		}
		return fmt.Sprintf("FAIL (%d/%d, %s)%s", passed, total, cancelledSeg, suffix)
	}
	frac := fmt.Sprintf("%d/%d", passed, total)
	if skipCount > 0 {
		frac = fmt.Sprintf("%d/%d, %d t.Skip", passed, total, skipCount)
	}
	if passed == total && !forceFail {
		token := fmt.Sprintf("PASS (%s)", frac)
		if style.enabled {
			return style.green(token) + suffix
		}
		return token + suffix
	}
	token := fmt.Sprintf("FAIL (%s)", frac)
	if style.enabled {
		return style.red(token) + suffix
	}
	return token + suffix
}

// PrintSkippedSummary writes a compact skip report grouped by label set.
// Paths and explanations appear only when verbose is true.
// Prefer PrintSkippedSummaryTo when a harness injects opts.Stdout.
func PrintSkippedSummary(skipped []core.SkippedCase, verbose bool) {
	PrintSkippedSummaryTo(os.Stdout, skipped, verbose)
}

// PrintSkippedSummaryTo is like PrintSkippedSummary but writes to w
// (nil means os.Stdout). Parallel-safe harnesses pass opts.Stdout.
func PrintSkippedSummaryTo(w io.Writer, skipped []core.SkippedCase, verbose bool) {
	if w == nil {
		w = os.Stdout
	}
	s := FormatSkippedSummary(skipped, verbose)
	if s == "" {
		return
	}
	fmt.Fprint(w, s)
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

// PrintResultSummary prints the end-of-run summary for a fully successful
// overall invocation (no prepare/run errors outside the pass counts).
func PrintResultSummary(opts core.Options, stats TestRunStats) {
	PrintResultSummaryOverall(opts, stats, true)
}

// PrintResultSummaryOverall prints PASS/FAIL. When overallOK is false, always
// uses FAIL even if passed==total (partial multi-tree: survivors passed but
// another tree failed prepare, or workspace error with odd counts).
//
// On timeout (stats.TimedOut) with cancelled > 0, the fraction uses Planned
// as denom and appends ", N cancelled" (orange when color on).
//
// On BuildFailed, prints FAIL (build failed; 0/planned executed) even when
// Total==0 (no suite leaves ran).
func PrintResultSummaryOverall(opts core.Options, stats TestRunStats, overallOK bool) {
	cancelled := timeoutCancelled(stats)
	// Timeout with all leaves cancelled may leave Total==0 (actual_run=0);
	// still print FAIL (0/planned, N cancelled). BuildFailed likewise.
	if stats.Total == 0 && stats.SkipCount == 0 && cancelled == 0 && !stats.GoTestBypassed && !stats.BuildFailed {
		return
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	style := newColorStyle(opts.Color, stdout)
	if stats.GoTestBypassed {
		planned := stats.Total
		if stats.Planned > 0 {
			planned = stats.Planned
		}
		fmt.Fprintln(stdout, formatBypassResultSummary(style, planned, stats.Elapsed))
		return
	}
	if stats.BuildFailed {
		planned := stats.Planned
		if planned <= 0 {
			planned = stats.Total
		}
		fmt.Fprintln(stdout, formatBuildFailedResultSummary(style, planned, stats.Elapsed))
		return
	}
	passed := stats.Passed
	total := stats.Total
	skipCount := stats.SkipCount
	if cancelled > 0 {
		if stats.Planned > 0 {
			total = stats.Planned
		}
		// v1: omit t.Skip phrase on timeout FAIL line.
		skipCount = 0
		overallOK = false
	}
	fmt.Fprintln(stdout, formatResultSummary(style, passed, total, stats.Elapsed, !overallOK, skipCount, cancelled))
}

// formatBypassResultSummary is the honest end line when go test was skipped.
func formatBypassResultSummary(style colorStyle, planned int, elapsed time.Duration) string {
	token := fmt.Sprintf("BYPASS (%d planned, 0 executed, go test bypassed)", planned)
	suffix := fmt.Sprintf(" in %s", formatDisplayDuration(elapsed))
	if style.enabled {
		// Neutral accent (same as muted info) — not green PASS.
		return token + suffix
	}
	return token + suffix
}

// formatBuildFailedResultSummary is the honest end line when the suite package
// failed to compile (no leaf subtests ran). Planned is discovery leaf count.
func formatBuildFailedResultSummary(style colorStyle, planned int, elapsed time.Duration) string {
	token := fmt.Sprintf("FAIL (build failed; 0/%d executed)", planned)
	suffix := fmt.Sprintf(" in %s", formatDisplayDuration(elapsed))
	if style.enabled {
		return style.red("FAIL") + fmt.Sprintf(" (build failed; 0/%d executed)", planned) + suffix
	}
	return token + suffix
}

func SkippedDisplayPath(doctestRoot, leafPath string) string {
	short := pathfmt.Short(doctestRoot)
	if leafPath == "" {
		return short
	}
	return short + "/" + filepath.ToSlash(leafPath)
}