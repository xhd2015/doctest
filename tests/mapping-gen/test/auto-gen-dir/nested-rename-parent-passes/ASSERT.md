---
label: heavy
---

## Expected
- Phase 1: `doctest test parent/nested` fails because `verbose_leaf` is a stale always-failing package.
- Phase 2: renaming `nested` to `nested-renamed` leaves stale generated tests under the old cache path.
- Phase 3: `doctest test parent` passes and only runs the parent-scoped leaf (`parent_ok`), ignoring stale nested cache packages.

## Exit Code
- Exit code 0.

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected parent run exit 0, got %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	if nestedRenameState.NestedFailResp == nil {
		t.Fatal("phase 1 nested run was not recorded")
	}
	if nestedRenameState.NestedFailResp.ExitCode == 0 {
		t.Fatalf("phase 1 expected nested verbose leaf to fail, got exit 0\nstderr:\n%s",
			nestedRenameState.NestedFailResp.Stderr)
	}
	phase1Out := nestedRenameState.NestedFailResp.Stdout + nestedRenameState.NestedFailResp.Stderr
	if !strings.Contains(phase1Out, "stale nested verbose_leaf package") {
		t.Fatalf("phase 1 expected stale leaf failure, output:\n%s", phase1Out)
	}

	if nestedRenameState.StaleCachePath == "" {
		t.Fatal("stale cache path not recorded")
	}
	if _, err := os.Stat(nestedRenameState.StaleCachePath); err != nil {
		t.Fatalf("expected stale generated package to remain at %s: %v", nestedRenameState.StaleCachePath, err)
	}

	if !strings.Contains(resp.Stderr, "1 tests") {
		t.Fatalf("parent run should discover only parent_ok, stderr:\n%s", resp.Stderr)
	}
	if strings.Contains(resp.Stdout, "verbose_leaf") || strings.Contains(resp.Stderr, "verbose_leaf") {
		t.Fatalf("parent run must not execute stale nested verbose_leaf package\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "PASS") {
		t.Fatalf("expected PASS in parent stdout:\n%s", resp.Stdout)
	}
}
```