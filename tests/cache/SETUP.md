# Scenario

**Feature**: inspect and clean durable doctest caches via in-process CLI

```
# L2: injectable CacheHome + cli.RunWithWriters
Harness -> t.TempDir CacheHome (+ optional outside overrides)
  -> seed CacheHome/doctest/<bucket>/...
  -> cli.RunWithWriters(["cache", flags...])
  -> stdout info | [dry-run] would remove | Removed
  -> filesystem side effects under req paths only
```

## Preconditions

- Nested root: does not inherit `tests/` binary `Run` or `testbin.Ensure`.
- All leaves are L2 in-process via `cli.RunWithWriters` + injectable roots on
  `Request`.
- Each leaf isolates its own `CacheHome` under `t.TempDir()`; never touches the
  developer's real user cache.
- No `os.Setenv` / `t.Setenv` / `os.Chdir` / `t.Chdir` in harness.
- Completeness: help, info, clean, flags, overrides (see DOCTEST.md tree).

## Steps

1. Root Setup is a no-op (helpers only).
2. Leaf Setup calls `ensureCacheHome`, optionally seeds buckets / overrides,
   sets `req.Args`.
3. `Run` calls `cli.RunWithWriters` (via `runWithInjectedCache`) and fills
   `Response`.

## Context

- `Request` / `Response` / `Run` are defined only in this tree's `DOCTEST.md`.
- Helpers below are tree-wide: root isolation, bucket seeding, path checks.
- **Layer**: L2 in-process CLI for all leaves.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// In-process only: no binary, no shared mutable process state.
	_ = d
	_ = req
	return nil
}

// ensureCacheHome sets req.CacheHome to t.TempDir() and DoctestRoot to
// CacheHome/doctest when CacheHome is empty. Does not create DoctestRoot
// (empty-root leaves leave it missing).
func ensureCacheHome(t *testing.T, req *Request) {
	t.Helper()
	if req.CacheHome == "" {
		req.CacheHome = t.TempDir()
	}
	if req.DoctestRoot == "" {
		req.DoctestRoot = filepath.Join(req.CacheHome, "doctest")
	}
}

// seedBucket writes a file under DoctestRoot/<name>/payload with the given
// content so the bucket has a known approximate size. Records name on
// req.SeededBuckets when not already listed.
func seedBucket(t *testing.T, req *Request, name string, content []byte) {
	t.Helper()
	ensureCacheHome(t, req)
	dir := filepath.Join(req.DoctestRoot, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir bucket %s: %v", name, err)
	}
	path := filepath.Join(dir, "payload")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write bucket %s: %v", name, err)
	}
	for _, b := range req.SeededBuckets {
		if b == name {
			return
		}
	}
	req.SeededBuckets = append(req.SeededBuckets, name)
}

// seedBytes returns a byte slice of length n (filled with 'x').
func seedBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'x'
	}
	return b
}

// pathExists reports whether path exists (file or directory).
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// combinedOut joins stdout and stderr for flexible message asserts.
func combinedOut(resp *Response) string {
	return resp.Stdout + "\n" + resp.Stderr
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

// mustContain fails if hay does not contain needle.
func mustContain(t *testing.T, hay, needle, label string) {
	t.Helper()
	if !strings.Contains(hay, needle) {
		t.Fatalf("%s: missing %q in:\n%s", label, needle, hay)
	}
}

// hasHumanSizeUnit reports whether s contains a human size unit token (B/K/M/G).
func hasHumanSizeUnit(s string) bool {
	// Accept 0B, 1.2K, 4.0K, 1M, 1.5G, etc.
	for _, u := range []string{"B", "K", "M", "G"} {
		if strings.Contains(s, u) {
			return true
		}
	}
	return false
}
```
