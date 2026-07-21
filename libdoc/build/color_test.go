package build

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

func TestFormatSkippedSummaryCompact(t *testing.T) {
	skipped := []core.SkippedCase{
		{DisplayPath: "a/heavy1", Labels: []string{"heavy"}},
		{DisplayPath: "a/heavy2", Labels: []string{"heavy"}},
		{DisplayPath: "b/slow", Labels: []string{"slow"}, Explanation: "takes time"},
		{DisplayPath: "c/both", Labels: []string{"slow", "heavy"}},
	}
	got := FormatSkippedSummary(skipped, false)
	if !strings.Contains(got, "skipped 4 labeled (discovery;") {
		t.Fatalf("header:\n%s", got)
	}
	if !strings.Contains(got, "heavy") || !strings.Contains(got, "2") {
		t.Fatalf("expected heavy bucket count 2:\n%s", got)
	}
	if !strings.Contains(got, "heavy,slow") {
		t.Fatalf("expected sorted multi-label key heavy,slow:\n%s", got)
	}
	if strings.Contains(got, "a/heavy1") {
		t.Fatalf("compact mode must not list paths:\n%s", got)
	}
	if !strings.Contains(got, "(use -v to list paths)") {
		t.Fatalf("expected -v hint:\n%s", got)
	}
	// Verbose lists paths + explanation.
	v := FormatSkippedSummary(skipped, true)
	if !strings.Contains(v, "a/heavy1") || !strings.Contains(v, "explanation: takes time") {
		t.Fatalf("verbose:\n%s", v)
	}
	if strings.Contains(v, "(use -v to list paths)") {
		t.Fatalf("verbose should not show -v hint:\n%s", v)
	}
}

func TestFormatSkippedSummaryLabelFilterHeader(t *testing.T) {
	skipped := []core.SkippedCase{
		{DisplayPath: "x", Labels: []string{"slow"}, Reason: "label filter"},
		{DisplayPath: "y", Labels: nil, Reason: "label filter"},
	}
	got := FormatSkippedSummary(skipped, false)
	if !strings.Contains(got, "skipped 2 (label filter;") {
		t.Fatalf("filter header:\n%s", got)
	}
	if !strings.Contains(got, "(unlabeled)") {
		t.Fatalf("unlabeled bucket:\n%s", got)
	}
}

func TestFormatDisplayDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{
			name: "sub-second nanoseconds to integer ms",
			d:    949802583 * time.Nanosecond,
			want: "949ms",
		},
		{
			name: "just over one second to two decimal places",
			d:    1366963417 * time.Nanosecond,
			want: "1.37s",
		},
		{
			name: "microseconds unchanged integer",
			d:    500 * time.Microsecond,
			want: "500µs",
		},
		{
			name: "integer milliseconds",
			d:    42 * time.Millisecond,
			want: "42ms",
		},
		{
			name: "exact one second",
			d:    time.Second,
			want: "1s",
		},
		{
			name: "sub-millisecond nanoseconds to integer ms",
			d:    1 * time.Millisecond,
			want: "1ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDisplayDuration(tt.d)
			if got != tt.want {
				t.Fatalf("formatDisplayDuration(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestFormatResultSummaryForceFail(t *testing.T) {
	style := colorStyle{enabled: false}
	elapsed := 2 * time.Second

	pass := formatResultSummary(style, 10, 10, elapsed, false, 0)
	if !strings.HasPrefix(pass, "PASS (10/10)") {
		t.Fatalf("expected PASS when all ok, got %q", pass)
	}
	// Survivors all passed but another tree failed prepare: must not look green.
	forced := formatResultSummary(style, 10, 10, elapsed, true, 0)
	if !strings.HasPrefix(forced, "FAIL (10/10)") {
		t.Fatalf("expected FAIL when forceFail, got %q", forced)
	}
	partial := formatResultSummary(style, 8, 10, elapsed, false, 0)
	if !strings.HasPrefix(partial, "FAIL (8/10)") {
		t.Fatalf("expected FAIL on partial pass, got %q", partial)
	}
}

func TestFormatResultSummaryRuntimeSkip(t *testing.T) {
	style := colorStyle{enabled: false}
	elapsed := time.Second

	// 1 pass + 1 t.Skip → succeeded/actual_run with t.Skip suffix.
	passSkip := formatResultSummary(style, 1, 1, elapsed, false, 1)
	if !strings.HasPrefix(passSkip, "PASS (1/1, 1 t.Skip) in ") {
		t.Fatalf("pass+skip: got %q", passSkip)
	}
	// 0 pass + 1 fail + 1 t.Skip.
	failSkip := formatResultSummary(style, 0, 1, elapsed, false, 1)
	if !strings.HasPrefix(failSkip, "FAIL (0/1, 1 t.Skip) in ") {
		t.Fatalf("fail+skip: got %q", failSkip)
	}
	// N=0 must keep legacy form (no t.Skip text).
	noSkip := formatResultSummary(style, 2, 2, elapsed, false, 0)
	if strings.Contains(noSkip, "t.Skip") || !strings.HasPrefix(noSkip, "PASS (2/2) in ") {
		t.Fatalf("zero skip must be legacy PASS (2/2), got %q", noSkip)
	}
}

func TestResolveColorMode(t *testing.T) {
	t.Run("always and never unchanged", func(t *testing.T) {
		var buf bytes.Buffer
		if got := ResolveColorMode(core.ColorAlways, &buf); got != core.ColorAlways {
			t.Fatalf("Always on buffer: got %v", got)
		}
		if got := ResolveColorMode(core.ColorNever, os.Stdout); got != core.ColorNever {
			t.Fatalf("Never on stdout: got %v", got)
		}
	})

	t.Run("auto on non-file is never", func(t *testing.T) {
		var buf bytes.Buffer
		if got := ResolveColorMode(core.ColorAuto, &buf); got != core.ColorNever {
			t.Fatalf("Auto on buffer: got %v, want Never", got)
		}
	})

	t.Run("auto on pipe is never", func(t *testing.T) {
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()
		if got := ResolveColorMode(core.ColorAuto, w); got != core.ColorNever {
			t.Fatalf("Auto on pipe: got %v, want Never", got)
		}
	})
}

// Regression: parallel ./... buffers progress into bytes.Buffer. ColorAuto
// against that buffer would always disable ANSI. The CLI resolves Auto against
// the real stdout first; after Always, writing into a buffer still emits color.
func TestColorAfterResolveIntoBuffer(t *testing.T) {
	var buf bytes.Buffer

	// Without resolve: Auto + buffer → plain.
	plain := newColorStyle(core.ColorAuto, &buf)
	if plain.enabled {
		t.Fatal("ColorAuto against buffer must disable color")
	}
	sumPlain := formatSummary(plain, 1, 1, 0, 0, time.Millisecond)
	if strings.Contains(sumPlain, "\x1b[") {
		t.Fatalf("expected plain summary, got %q", sumPlain)
	}

	// After resolve (as runner.Test does against user-facing stdout): Always
	// into the same buffer shape → colored Pass segment.
	resolved := ResolveColorMode(core.ColorAlways, &buf) // explicit Always after TTY resolve
	style := newColorStyle(resolved, &buf)
	if !style.enabled {
		t.Fatal("ColorAlways against buffer must enable color")
	}
	sum := formatSummary(style, 1, 1, 0, 0, time.Millisecond)
	if !strings.Contains(sum, ansiGreen+"1 Pass"+ansiReset) {
		t.Fatalf("expected green 1 Pass in summary, got %q", sum)
	}
	if !strings.Contains(sum, ansiGray+"0 Fail"+ansiReset) {
		t.Fatalf("expected gray 0 Fail in summary, got %q", sum)
	}
}