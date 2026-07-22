---
label: heavy
---

## Expected

- Non-zero exit.
- Stderr reports a clear error refusing default / warm `mapping-gen` as cold gen-dir.
- Warm mapping-gen content is **not** wiped (marker still present).

## Errors

- Refuse `--gen-dir` equal to warm mapping-gen home.

## Exit Code

- non-zero

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit when --gen-dir equals warm mapping-gen, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}

	combined := strings.ToLower(resp.Stderr + "\n" + resp.Stdout)
	// Prefer a product error about mapping-gen / refuse — not only a generic panic.
	mentionsWarm := strings.Contains(combined, "mapping-gen") ||
		strings.Contains(combined, "warm") ||
		strings.Contains(combined, "cold-cache") ||
		strings.Contains(combined, "cold cache")
	refuses := strings.Contains(combined, "error") ||
		strings.Contains(combined, "refus") ||
		strings.Contains(combined, "not allow") ||
		strings.Contains(combined, "cannot") ||
		strings.Contains(combined, "must not") ||
		strings.Contains(combined, "invalid")
	if !mentionsWarm || !refuses {
		t.Fatalf("expected clear Error refusing warm mapping-gen gen-dir, got:\nstderr:\n%s\nstdout:\n%s", resp.Stderr, resp.Stdout)
	}

	if st.Marker == "" {
		t.Fatal("st.Marker not set")
	}
	if _, statErr := os.Stat(st.Marker); statErr != nil {
		t.Fatalf("warm mapping-gen was wiped or marker missing (must not wipe on reject): %s: %v\nstderr:\n%s", st.Marker, statErr, resp.Stderr)
	}
}
```
