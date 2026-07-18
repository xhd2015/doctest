package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

type TestRunStats struct {
	Passed         int
	Total          int
	Elapsed        time.Duration
	NoTestsChanged bool
	Skipped        []core.SkippedCase
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

func PrintSkippedSummary(skipped []core.SkippedCase) {
	if len(skipped) == 0 {
		return
	}
	fmt.Printf("SKIPPED %d TESTS\n", len(skipped))
	for _, s := range skipped {
		fmt.Printf("  %s\n", s.DisplayPath)
		if len(s.Labels) > 0 {
			fmt.Printf("    label: %s\n", strings.Join(s.Labels, ", "))
		}
		if s.Explanation != "" {
			fmt.Printf("    explanation: %s\n", s.Explanation)
		}
		if s.Reason != "" {
			fmt.Printf("    reason: %s\n", s.Reason)
		}
	}
	fmt.Println()
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