# Scenario

**Feature**: inventory doctest roots via in-process `cli.RunWithWriters` (no product binary)

```
# fixture trees under t.TempDir
Harness -> write DOCTEST.md + ASSERT leaves
  -> cli.RunWithWriters(["list", absPatterns...])
  -> stdout body lines + summary; stderr soft/errors
```

## Preconditions

- Nested root: does not inherit `tests/` binary `Run` or `testbin.Ensure`.
- All leaves are L2 in-process via `cli.RunWithWriters` + injected writers.
- Fixture trees live under `t.TempDir()`; patterns are absolute (or
  `absBase + "/..."`) so process cwd need not be the fixture (cwd undetermined).
- No `os.Setenv` / `t.Setenv` / `os.Chdir` / `t.Chdir` in harness.
- Color leaves use `--color` / `--no-color` only (no process-global `NO_COLOR`).
- Completeness: help, discovery, inventory, summary, color (see DOCTEST.md tree).

## Steps

1. Root Setup is a no-op (helpers only).
2. Leaf Setup writes fixture trees and sets `req.Args` / `req.Roots` / `req.FixtureDir`.
3. `Run` calls `cli.RunWithWriters` and fills `Response`.

## Context

- `Request` / `Response` / `Run` are defined only in this tree's `DOCTEST.md`.
- Helpers below are tree-wide: minimal fixture writers + stdout format checkers.
- **Layer**: L2 in-process CLI for all leaves.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// In-process only: no binary, no shared mutable process state.
	_ = d
	_ = req
	return nil
}

// fixtureScenarioSETUP is minimal SETUP.md content for inventory fixtures.
// Fence markers use hex escapes so this root SETUP does not close its own go fence.
const fixtureScenarioSETUP = "# Scenario\n\n**Feature**: minimal fixture\n\n\x60\x60\x60\n# fixture pipeline\nsystem -> run\n\x60\x60\x60\n"

// writeRootDOCTEST writes a minimal valid DOCTEST.md (Request/Response/Run) at dir.
func writeRootDOCTEST(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteFile(t, dir, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
	testtree.WriteFile(t, dir, "SETUP.md", fixtureScenarioSETUP)
}

// writeLeafASSERT writes root/rel/ASSERT.md with optional frontmatter labels.
// labels is a comma-separated label field (e.g. "e2e, slow") or empty for unlabeled.
func writeLeafASSERT(t *testing.T, root, rel, labels string) {
	t.Helper()
	dir := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	testtree.WriteFile(t, dir, "SETUP.md", fixtureScenarioSETUP)
	body := "## Expected\n- fixture leaf\n"
	if labels != "" {
		body = "---\nlabel: " + labels + "\n---\n\n" + body
	}
	testtree.WriteFile(t, dir, "ASSERT.md", body)
}

// writeLabeledLeaves creates unlabeled and labeled leaves under root.
// Each entry is "rel" or "rel|labelField" (labelField e.g. "e2e" or "e2e, slow").
func writeLabeledLeaves(t *testing.T, root string, specs []string) {
	t.Helper()
	writeRootDOCTEST(t, root)
	for _, spec := range specs {
		rel, labels, _ := strings.Cut(spec, "|")
		writeLeafASSERT(t, root, rel, labels)
	}
}

// listArgs builds ["list", extraFlags..., patterns...].
func listArgs(flags []string, patterns ...string) []string {
	args := []string{"list"}
	args = append(args, flags...)
	args = append(args, patterns...)
	return args
}

// bodyLineRE matches one uncolored body line:
// path \t leaves \t L2:L3=a:b [ (p2%/p3%)] \t labelDist
var bodyLineRE = regexp.MustCompile(`^([^\t]+)\t(\d+)\tL2:L3=(\d+):(\d+)(?: \(([0-9.]+)%/([0-9.]+)%\))?\t(.+)$`)

// summaryTotalsRE matches the totals summary line (uncolored).
var summaryTotalsRE = regexp.MustCompile(`^roots=(\d+)  leaves=(\d+)  L2:L3=(\d+):(\d+)(?:  \(L2 ([0-9.]+)% / L3 ([0-9.]+)%\))?$`)

type parsedBody struct {
	Path     string
	Leaves   int
	L2, L3   int
	HasPct   bool
	P2, P3   string // percent strings without %
	LabelDist string
}

type parsedReport struct {
	Body    []parsedBody
	HasSep  bool // blank line + ---
	Totals  string
	Labels  string
	Raw     string
}

func parseListReport(t *testing.T, stdout string) parsedReport {
	t.Helper()
	rep := parsedReport{Raw: stdout}
	// Require trailing newline on non-empty successful reports (POSIX CLI).
	if stdout == "" {
		return rep
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout must end with trailing newline:\n%q", stdout)
	}
	// Strip final newline for split; keep empty last line semantics via Split after TrimSuffix.
	trim := strings.TrimSuffix(stdout, "\n")
	lines := strings.Split(trim, "\n")
	// Find separator "---" after a blank line.
	sepIdx := -1
	for i, ln := range lines {
		if ln == "---" && i > 0 && lines[i-1] == "" {
			sepIdx = i
			break
		}
	}
	var bodyLines []string
	if sepIdx < 0 {
		// No summary yet (or empty) — treat all non-empty as body (caller checks).
		for _, ln := range lines {
			if ln != "" {
				bodyLines = append(bodyLines, ln)
			}
		}
	} else {
		rep.HasSep = true
		for _, ln := range lines[:sepIdx-1] { // exclude the blank before ---
			if ln == "" {
				continue
			}
			bodyLines = append(bodyLines, ln)
		}
		// After ---: totals then labels
		rest := lines[sepIdx+1:]
		if len(rest) < 2 {
			t.Fatalf("summary after --- needs totals + labels lines, got %v\nfull:\n%s", rest, stdout)
		}
		rep.Totals = rest[0]
		rep.Labels = rest[1]
		if len(rest) > 2 {
			t.Fatalf("extra lines after summary labels: %v\nfull:\n%s", rest[2:], stdout)
		}
	}
	for _, ln := range bodyLines {
		m := bodyLineRE.FindStringSubmatch(ln)
		if m == nil {
			t.Fatalf("body line does not match format:\n%q\nfull stdout:\n%s", ln, stdout)
		}
		leaves, _ := strconv.Atoi(m[2])
		l2, _ := strconv.Atoi(m[3])
		l3, _ := strconv.Atoi(m[4])
		pb := parsedBody{
			Path:      m[1],
			Leaves:    leaves,
			L2:        l2,
			L3:        l3,
			HasPct:    m[5] != "",
			P2:        m[5],
			P3:        m[6],
			LabelDist: m[7],
		}
		if pb.L2+pb.L3 != pb.Leaves {
			t.Fatalf("L2+L3 != leaves on line %q", ln)
		}
		if pb.Leaves == 0 && pb.HasPct {
			t.Fatalf("percent group must be omitted when leaves==0: %q", ln)
		}
		if pb.Leaves > 0 && !pb.HasPct {
			t.Fatalf("percent group required when leaves>0: %q", ln)
		}
		rep.Body = append(rep.Body, pb)
	}
	return rep
}

func requireOK(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%q stdout=%q err=%v", resp.ExitCode, resp.Stderr, resp.Stdout, resp.Err)
	}
}

func requireFail(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit; stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
}

func requireSortedPaths(t *testing.T, body []parsedBody) {
	t.Helper()
	paths := make([]string, len(body))
	for i, b := range body {
		paths[i] = b.Path
	}
	sorted := append([]string(nil), paths...)
	sort.Strings(sorted)
	for i := range paths {
		if paths[i] != sorted[i] {
			t.Fatalf("body paths not sorted:\n got %v\nwant %v", paths, sorted)
		}
	}
}

func findBody(t *testing.T, body []parsedBody, absPath string) parsedBody {
	t.Helper()
	// Product may print path as-given or cleaned; compare with Abs when needed.
	want, err := filepath.Abs(absPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range body {
		got, gerr := filepath.Abs(b.Path)
		if gerr != nil {
			got = b.Path
		}
		if got == want || b.Path == absPath || b.Path == want {
			return b
		}
	}
	var got []string
	for _, b := range body {
		got = append(got, b.Path)
	}
	t.Fatalf("path %q not in body paths %v", absPath, got)
	return parsedBody{}
}

func requireLabelDistHas(t *testing.T, dist, key string, count int) {
	t.Helper()
	want := fmt.Sprintf("%s=%d", key, count)
	// labelDist is space-separated name=count tokens
	for _, tok := range strings.Fields(dist) {
		if tok == want {
			return
		}
	}
	t.Fatalf("labelDist %q missing %s", dist, want)
}

func requireNoANSI(t *testing.T, s, which string) {
	t.Helper()
	if strings.Contains(s, "\x1b") {
		t.Fatalf("%s must not contain ANSI ESC:\n%q", which, s)
	}
}

func requireGrayMeta(t *testing.T, stdout string) {
	t.Helper()
	// Gray SGR used by build color helpers: \x1b[90m
	if !strings.Contains(stdout, "\x1b[90m") {
		t.Fatalf("expected gray SGR (\\x1b[90m) on meta when --color:\n%q", stdout)
	}
	// Path field itself should remain plain: first field before tab should not open with ESC.
	// With colored meta, line may still start with plain path.
	for _, ln := range strings.Split(strings.TrimSuffix(stdout, "\n"), "\n") {
		if ln == "" || ln == "---" || strings.HasPrefix(ln, "roots=") || strings.HasPrefix(ln, "labels:") {
			continue
		}
		// body line: path is first tab field and must not start with ESC
		tab := strings.IndexByte(ln, '\t')
		if tab < 0 {
			continue
		}
		pathField := ln[:tab]
		if strings.Contains(pathField, "\x1b") {
			t.Fatalf("path field must be plain (no ANSI): %q", pathField)
		}
	}
}

func fmtPct(n, total int) string {
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("%.1f", 100.0*float64(n)/float64(total))
}
```
