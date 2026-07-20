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