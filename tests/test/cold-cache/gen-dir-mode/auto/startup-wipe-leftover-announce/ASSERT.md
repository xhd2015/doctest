---
label: heavy
---

## Expected

- Command succeeds (exit 0).
- Startup wipe: pre-seeded `marker-before` under cold home is gone.
- Leftover: cold home still exists and contains generated Go content (not deleted on finish).
- Announce: stderr describes cold-cache mode (keyword and/or cold gen path; GOCACHE isolation and count mentioned).

## Side Effects

- `$DOCTEST_CACHE_HOME/doctest/mapping-gen-cold` recreated and populated; not removed after run.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	if st.Marker == "" || st.ColdHome == "" {
		t.Fatal("st.Marker / st.ColdHome not set by setup")
	}
	if _, statErr := os.Stat(st.Marker); statErr == nil {
		t.Fatalf("startup wipe failed: marker still exists at %s", st.Marker)
	} else if !os.IsNotExist(statErr) {
		t.Fatalf("stat marker: %v", statErr)
	}

	fi, statErr := os.Stat(st.ColdHome)
	if statErr != nil {
		t.Fatalf("cold home missing after finish (should leftover): %s: %v\nstderr:\n%s", st.ColdHome, statErr, resp.Stderr)
	}
	if !fi.IsDir() {
		t.Fatalf("cold home is not a directory: %s", st.ColdHome)
	}
	if !dirHasGoFiles(st.ColdHome) {
		// Also accept non-empty tree if layout uses only go.mod at root of cold home.
		entries, readErr := os.ReadDir(st.ColdHome)
		if readErr != nil || len(entries) == 0 {
			t.Fatalf("cold home empty after run (expected generated leftover): %s\nstderr:\n%s", st.ColdHome, resp.Stderr)
		}
	}

	combined := resp.Stderr + "\n" + resp.Stdout
	lower := strings.ToLower(combined)
	announcesCold := strings.Contains(lower, "cold-cache") ||
		strings.Contains(lower, "cold cache") ||
		strings.Contains(combined, "mapping-gen-cold") ||
		strings.Contains(combined, st.ColdHome)
	if !announcesCold {
		t.Fatalf("expected cold-cache announcement (cold-cache keyword and/or cold gen path) on stderr:\n%s", resp.Stderr)
	}
	// Announcement should cover isolated GOCACHE and count when present in the product line.
	if !strings.Contains(lower, "gocache") && !strings.Contains(lower, "go cache") {
		// Soft preference: still require gen/cold signal above; GOCACHE word is part of locked announce.
		t.Fatalf("expected cold-cache announcement to mention GOCACHE isolation:\n%s", resp.Stderr)
	}
	if !strings.Contains(lower, "count") {
		t.Fatalf("expected cold-cache announcement to mention count:\n%s", resp.Stderr)
	}

	// Sanity: marker name must not reappear as leftover content root file.
	if _, err := os.Stat(filepath.Join(st.ColdHome, "marker-before")); err == nil {
		t.Fatalf("marker-before reappeared under cold home")
	}
}
```
