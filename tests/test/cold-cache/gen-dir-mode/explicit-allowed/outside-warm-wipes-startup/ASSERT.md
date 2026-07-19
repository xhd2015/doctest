## Expected

- Command succeeds (exit 0).
- Explicit `--gen-dir` was wiped on startup (`marker-before` gone).
- Explicit gen dir still exists after finish with generated content (leftover).

## Side Effects

- Only the explicit gen dir is wiped/recreated — not assert-mod / session-mod / selftest-bin.

## Exit Code

- 0

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if st.GenDir == "" || st.Marker == "" {
		t.Fatal("st.GenDir / st.Marker not set")
	}
	if _, statErr := os.Stat(st.Marker); statErr == nil {
		t.Fatalf("startup wipe failed: marker still exists at %s", st.Marker)
	} else if !os.IsNotExist(statErr) {
		t.Fatalf("stat marker: %v", statErr)
	}

	fi, statErr := os.Stat(st.GenDir)
	if statErr != nil {
		t.Fatalf("explicit gen dir missing after finish: %s: %v\nstderr:\n%s", st.GenDir, statErr, resp.Stderr)
	}
	if !fi.IsDir() {
		t.Fatalf("explicit gen dir is not a directory: %s", st.GenDir)
	}
	if !dirHasGoFiles(st.GenDir) {
		entries, readErr := os.ReadDir(st.GenDir)
		if readErr != nil || len(entries) == 0 {
			t.Fatalf("explicit gen dir empty after run: %s\nstderr:\n%s", st.GenDir, resp.Stderr)
		}
	}
}
```
