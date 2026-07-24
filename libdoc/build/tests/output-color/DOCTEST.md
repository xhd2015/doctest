# Colored `doctest test` Output

## Version
0.0.2


Tests that non-verbose `build.Test` progress output optionally emits ANSI color
on dot progress and the summary line, controlled by `core.Options.Color`
(`ColorAuto`, `ColorAlways`, `ColorNever`).

## DSN (Domain Specific Notion)

### Participants
- **`build.Test`** — runs generated Go tests, prints one dot per package on
  stdout, then a summary line `(N Run, N Pass, N Fail, N Cached)`.
- **`ColorMode`** — `ColorAuto` (TTY detection), `ColorAlways` (`--color`),
  `ColorNever` (`--no-color`).
- **ANSI helpers** — wrap metric substrings and fail dots when color is enabled.

### Behaviors
- **Dot progress** — fail dots red when color on; pass and cached dots plain.
- **Summary line** — `N Pass` green when N>0; `N Fail` red when N>0, gray when
  N=0; `N Cached` gray always; `N Run` plain.
- **Color off** — plain dots and summary (auto on non-TTY pipe, or never).
- **CLI `./...`** — `runner.Test` resolves `ColorAuto` against the real stdout
  via `build.ResolveColorMode` *before* buffering each tree into a
  `bytes.Buffer` (parallel non-interleave). Otherwise Auto would always see a
  non-file writer and disable color even on a TTY.

## Decision Tree

```
output-color
├── color-disabled          [Color off]
│   ├── auto-on-pipe        ColorAuto + opts.Stdout buffer (non-TTY) → no ANSI
│   └── never-on-fail       ColorNever + failing leaf → no ANSI anywhere
└── color-enabled           [ColorAlways — deterministic]
    ├── force-color         1 pass → green `N Pass`, plain dot
    ├── all-pass            2 pass → plain dots; green Pass; gray Fail/Cached (0)
    ├── has-fail            1 pass + 1 fail → red fail dot; green Pass; red Fail
    └── cached-gray         warm cache → gray `N Cached` segment (N>0)
```

## Test Index

| Leaf | Description |
|------|-------------|
| `color-disabled/auto-on-pipe` | `ColorAuto` with `opts.Stdout` buffer (non-TTY) emits no escape sequences |
| `color-disabled/never-on-fail` | `ColorNever` leaves fail dot and summary uncolored |
| `color-enabled/force-color` | `ColorAlways` wraps `N Pass` in green and `0 Fail` in gray |
| `color-enabled/all-pass` | Two pass dots stay plain; Pass green; `0 Fail` and `0 Cached` gray |
| `color-enabled/has-fail` | Fail dot red; Pass green; Fail red; pass dot plain |
| `color-enabled/cached-gray` | After cache warmup, `N Cached` and `0 Fail` summary segments are gray |

## How to Run

```sh
doctest test ./libdoc/build/tests/output-color/...
doctest test ./tests/help/test-options
doctest test ./tests/test/color-flags/conflict
```

```go
import (
	"bytes"
	"path/filepath"
	"testing"
	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
)

type Request struct {
	Color		core.ColorMode
	Count		int
	PassCount	int
	FailCount	int
	WarmCache	bool
}
type Response struct {
	Output	string
	Dots	string
	Summary	string
	GenDir	string
	TestErr	error
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	subRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gendir")

	createSubTree(t, subRoot, req.PassCount, req.FailCount)

	// Parallel-safe capture: write progress into opts.Stdout (bytes.Buffer is
	// non-TTY, so ColorAuto resolves off). Never swap os.Stdout — that races
	// under t.Parallel() across leaves in the same package.
	var buf bytes.Buffer
	opts := core.Options{
		GenDir:		genDir,
		RemoveTemp:	false,
		Color:		req.Color,
		Count:		req.Count,
		Stdout:		&buf,
		// Keep PASS/FAIL end-line off real stdout during harness capture.
		SuppressResultSummary:	true,
	}

	if req.WarmCache {
		// Warmup must not pollute the measured buffer (and must not race
		// os.Stdout). Reuse the same opts but discard progress.
		var warmDiscard bytes.Buffer
		warmOpts := opts
		warmOpts.Stdout = &warmDiscard
		// First run may rewrite gen with -count=1 (no result cache store);
		// second run stores the entry the measured run should hit.
		if err := build.Test(subRoot, warmOpts); err != nil {
			t.Fatalf("cache warmup run 1 failed: %v", err)
		}
		if err := build.Test(subRoot, warmOpts); err != nil {
			t.Fatalf("cache warmup run 2 failed: %v", err)
		}
	}

	testErr := build.Test(subRoot, opts)

	output := buf.String()
	dots, summary := splitDotsAndSummary(output)

	return &Response{
		Output:		output,
		Dots:		dots,
		Summary:	summary,
		GenDir:		genDir,
		TestErr:	testErr,
	}, nil
}
```
